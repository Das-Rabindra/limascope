package main

import (
	"context"
	"crypto/tls"
	"embed"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amir20/dozzle/internal/agent"
	"github.com/amir20/dozzle/internal/auth"
	"github.com/amir20/dozzle/internal/docker"
	"github.com/amir20/dozzle/internal/k8s"
	"github.com/amir20/dozzle/internal/support/cli"
	docker_support "github.com/amir20/dozzle/internal/support/docker"
	k8s_support "github.com/amir20/dozzle/internal/support/k8s"
	"github.com/amir20/dozzle/internal/web"
	"github.com/rs/zerolog/log"
)

//go:embed all:dist
var content embed.FS

func main() {
	cli.ValidateEnvVars(cli.Args{}, cli.AgentCmd{})
	args, subcommand := cli.ParseArgs()

	if subcommand != nil {
		runnable, ok := subcommand.(cli.Runnable)
		if !ok {
			log.Fatal().Msg("Invalid command")
		}
		err := runnable.Run(args, content)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run command")
		}
		os.Exit(0)
	}

	if args.AuthProvider != "none" && args.AuthProvider != "forward-proxy" && args.AuthProvider != "simple" {
		log.Fatal().Str("provider", args.AuthProvider).Msg("Invalid auth provider")
	}

	log.Info().Msgf("Dozzle version %s", args.Version())

	var hostService web.HostService
	switch args.Mode {
	case "server":
		multiHostService := cli.CreateMultiHostService(content, args)
		if multiHostService.TotalClients() == 0 {
			log.Fatal().Msg("Could not connect to any Docker Engine")
		}
		log.Info().Int("clients", multiHostService.TotalClients()).Msg("Connected to Docker")
		hostService = multiHostService

	case "swarm":
		localClient, err := docker.NewLocalClient("")
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create docker client")
		}

		dummyCert := tls.Certificate{}
		agentManager := docker_support.NewRetriableClientManager(args.RemoteAgent, args.Timeout, dummyCert)
		manager := docker_support.NewSwarmClientManager(localClient, dummyCert, args.Timeout, agentManager, args.Filter)
		hostService = docker_support.NewMultiHostService(manager, args.Timeout)

		log.Info().Msg("Starting in swarm mode")
		listener, err := net.Listen("tcp", ":7007")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen")
		}

		server, err := agent.NewServer(localClient, dummyCert, args.Version(), args.Filter)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create agent")
		}

		go cli.StartEvent(args, "swarm", localClient, "")
		go func() {
			log.Info().Msgf("Dozzle agent version in swarm mode %s", args.Version())
			if err := server.Serve(listener); err != nil {
				log.Error().Err(err).Msg("failed to serve")
			}
		}()

	case "k8s":
		localClient, err := k8s.NewK8sClient(args.Namespace)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create k8s client")
		}

		clusterService, err := k8s_support.NewK8sClusterService(localClient, args.Timeout)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create k8s cluster service")
		}

		go cli.StartEvent(args, "k8s", localClient, "")
		hostService = clusterService

	default:
		log.Fatal().Str("mode", args.Mode).Msg("Invalid mode")
	}

	srv := createServer(args, hostService)
	go func() {
		log.Info().Msgf("Accepting connections on %s", args.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to listen")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	stop()
	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shut down")
	}
	log.Debug().Msg("shut down complete")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func createServer(args cli.Args, hostService web.HostService) *http.Server {
	_, dev := os.LookupEnv("DEV")

	var provider web.AuthProvider = web.NONE
	var authorizer web.Authorizer

	if args.AuthProvider == "forward-proxy" {
		provider = web.FORWARD_PROXY
		authorizer = auth.NewForwardProxyAuth(args.AuthHeaderUser, args.AuthHeaderEmail, args.AuthHeaderName, args.AuthHeaderFilter)
	} else if args.AuthProvider == "simple" {
		provider = web.SIMPLE

		userFilePath := "./data/users.yml"
		if !fileExists(userFilePath) {
			userFilePath = "./data/users.yaml"
			if !fileExists(userFilePath) {
				log.Fatal().Msg("No users.yaml or users.yml file found.")
			}
		}

		db, err := auth.ReadUsersFromFile(userFilePath)
		if err != nil {
			log.Fatal().Err(err).Msgf("Could not read users file: %s", userFilePath)
		}
		ttl := time.Duration(0)
		if args.AuthTTL != "session" {
			ttl, err = time.ParseDuration(args.AuthTTL)
			if err != nil {
				log.Fatal().Err(err).Msg("Could not parse auth ttl")
			}
		}
		authorizer = auth.NewSimpleAuth(db, ttl)
	}

	authTTL := time.Duration(0)
	if args.AuthTTL != "session" {
		ttl, err := time.ParseDuration(args.AuthTTL)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not parse auth ttl")
		}
		authTTL = ttl
	}

	config := web.Config{
		Addr:        args.Addr,
		Base:        args.Base,
		Version:     args.Version(),
		Hostname:    args.Hostname,
		NoAnalytics: args.NoAnalytics,
		Dev:         dev,
		Authorization: web.Authorization{
			Provider:   provider,
			Authorizer: authorizer,
			TTL:        authTTL,
		},
		EnableActions: args.EnableActions,
		EnableShell:   args.EnableShell,
		Labels:        args.Filter,
	}

	assets, err := fs.Sub(content, "dist")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not get sub filesystem")
	}

	if _, ok := os.LookupEnv("LIVE_FS"); ok {
		if dev {
			assets = os.DirFS("./public")
		} else {
			assets = os.DirFS("./dist")
		}
	}

	if !dev {
		if _, err := assets.Open(".vite/manifest.json"); err != nil {
			log.Fatal().Msg("manifest.json not found")
		}
		if _, err := assets.Open("index.html"); err != nil {
			log.Fatal().Msg("index.html not found")
		}
	}

	return web.CreateServer(hostService, assets, config)
}
