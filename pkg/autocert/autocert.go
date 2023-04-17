// Package autocert support auto cert
package autocert

import (
	"crypto/tls"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

var (
	once    sync.Once
	manager *Manager
)

// Manager is autocert Manager Proxy
type Manager struct {
	*autocert.Manager
	domain string
}

// New create a autocert Manager
func New(domain string) {
	once.Do(func() {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("./certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
		}
		manager = &Manager{
			m,
			domain,
		}
	})
}

// GetTLSCert get tls Certificate
func GetTLSCert() (*tls.Certificate, error) {
	cert, err := manager.TLSConfig().GetCertificate(ClientHelloInfo(manager.domain))
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// GetTLSConfig get tls Config
func GetTLSConfig() *tls.Config {
	if _, err := GetTLSCert(); err != nil {
		return nil
	}
	return manager.TLSConfig()
}

// ClientHelloInfo get tls ClientHelloInfo
func ClientHelloInfo(sni string) *tls.ClientHelloInfo {
	hello := &tls.ClientHelloInfo{
		ServerName:   sni,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
	}
	return hello
}
