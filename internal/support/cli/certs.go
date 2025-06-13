package cli

import (
	"crypto/tls"
)

// ReadCertificates returns an empty TLS certificate for development mode
func ReadCertificates(_ any) (tls.Certificate, error) {
	return tls.Certificate{}, nil
}
