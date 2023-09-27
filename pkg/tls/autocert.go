// Package tls support tls management
package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/acme/autocert"

	"github.com/beihai0xff/pudding/configs"
)

var (
	once    sync.Once
	manager *Manager
)

// ConfigType is a tls config type
type ConfigType int

const (
	// ConfigTypeServer is a tls config type for server
	ConfigTypeServer ConfigType = iota
	// ConfigTypeClient is a tls config type for client
	ConfigTypeClient ConfigType = iota
)

// Manager is a tls management struct
type Manager struct {
	autoCert        *autocert.Manager
	domain          string
	serverTLSConfig *tls.Config
	clientTLSConfig *tls.Config
}

// New create an tls Manager
func New(hostDomain string, conf *configs.TLS) {
	once.Do(func() {
		if conf.AutoCert {
			m := &autocert.Manager{
				Cache:      autocert.DirCache("./certs"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(hostDomain),
			}

			manager = &Manager{
				domain:   hostDomain,
				autoCert: m,
			}

			return
		}

		if err := readCertFromFile(conf); err != nil {
			panic(err)
		}
	})
}

func readCertFromFile(conf *configs.TLS) error {
	// Load CA cert
	ca, err := os.ReadFile(conf.CACert)
	if err != nil {
		return fmt.Errorf("ReadFile %s failed: %w", conf.CACert, err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(ca); !ok {
		return fmt.Errorf("AppendCertsFromPEM %s failed", conf.CACert)
	}

	// Load server cert and key
	serverCert, err := tls.LoadX509KeyPair(conf.ServerCert, conf.ServerKey)
	if err != nil {
		panic(err)
	}

	// Load client cert and key
	clientCert, err := tls.LoadX509KeyPair(conf.ClientCert, conf.ClientKey)
	if err != nil {
		panic(err)
	}

	manager = &Manager{
		serverTLSConfig: &tls.Config{
			ClientCAs:    caCertPool,
			Certificates: []tls.Certificate{serverCert},
			ClientAuth:   tls.VerifyClientCertIfGiven,
			MinVersion:   tls.VersionTLS13,
			NextProtos:   []string{"h2", "http/1.1"}, // enable HTTP/2
		},
		clientTLSConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		},
	}

	return nil
}

// GetTLSConfig get tls Config
func GetTLSConfig(t ConfigType) *tls.Config {
	if manager.autoCert != nil {
		return manager.autoCert.TLSConfig()
	}

	switch t {
	case ConfigTypeServer:
		return manager.serverTLSConfig
	case ConfigTypeClient:
		return manager.clientTLSConfig
	default:
		return nil
	}
}

// ClientHelloInfo get tls ClientHelloInfo
func ClientHelloInfo(sni string) *tls.ClientHelloInfo {
	hello := &tls.ClientHelloInfo{
		ServerName:   sni,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305, tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305},
	}

	return hello
}
