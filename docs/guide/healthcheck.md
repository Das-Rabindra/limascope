---
title: Healthcheck
---

# Enabling Healthcheck

Limascope has internal support for healthcheck using the `limascope healthcheck` command. It is not enabled by default as it adds extra CPU usage. To use `healthcheck`, you need to configure it. Below is an example that checks the health of Limascope every 3 seconds.

```yaml
services:
  limascope:
    image: Das-Rabindra/limascope:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - 8080:8080
    healthcheck:
      test: ["CMD", "/limascope", "healthcheck"]
      interval: 3s
      timeout: 30s
      retries: 5
      start_period: 30s
```

`limascope healthcheck` skips agents as they are not required for healthcheck. Agents can be configured to have their own [healthcheck](/guide/agent#setting-up-healthcheck).

> [!WARNING]
> The `healthcheck` command does not work with `--health-cmd` flag due to a bug in Docker. You need to use the `healthcheck` configuration in the `docker-compose.yml` file. See [Docker issue](https://github.com/docker/cli/issues/3719) for more information.
