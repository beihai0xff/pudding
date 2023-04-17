package autocert

import (
	"crypto/tls"
	"net"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

var (
	once    sync.Once
	manager *Manager
)

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

func GetCertificates() (*tls.Certificate, error) {
	cert, err := manager.TLSConfig().GetCertificate(ClientHelloInfo(manager.domain))
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func GetTLSConfig() *tls.Config {
	if _, err := GetCertificates(); err != nil {
		return nil
	}
	return manager.TLSConfig()
}

func GetListener() net.Listener {
	return manager.Listener()
}

func ClientHelloInfo(sni string) *tls.ClientHelloInfo {
	hello := &tls.ClientHelloInfo{
		ServerName:   sni,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
	}
	return hello
}
