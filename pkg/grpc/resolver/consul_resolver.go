// Package resolver provides a resolver interface for service discovery.
// consul_resolver.go provides a consul resolver.
package resolver

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/utils"
)

// consulResolver
type consulResolver struct {
	client *api.Client
}

var (
	consulOnce sync.Once
	consul     Resolver
)

// WithConsulResolver new
func WithConsulResolver(addr string) OptionResolver {
	consulOnce.Do(func() {
		cfg := api.DefaultConfig()
		cfg.Address = addr
		c, err := api.NewClient(cfg)
		if err != nil {
			log.Fatalf("failed to create consul client: %v", err)
		}
		consul = &consulResolver{client: c}
	})

	return func() Resolver {
		return consul
	}
}

// RegisterGRPC register gRPC service to consul
func (c *consulResolver) RegisterGRPC(serviceName, ip string, port int) (string, error) {
	// health check
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, port),
		GRPCUseTLS:                     true,
		TLSSkipVerify:                  true,
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}

	serviceID := fmt.Sprintf("%s-%s:%d", serviceName, ip, port)

	srv := &api.AgentServiceRegistration{
		ID:      serviceID,                                   // service unique ID
		Name:    "grpc." + serviceName,                       // service name
		Tags:    []string{utils.GetEnv(), "pudding", "gRPC"}, // service tags
		Address: ip,
		Port:    port,
		Check:   check,
	}

	if err := c.client.Agent().ServiceRegister(srv); err != nil {
		log.Errorf("register grpc service [%s] failed: %v", serviceID, err)
		return "", err
	}

	log.Infof("register grpc service [%s] successfully", serviceID)

	return serviceID, nil
}

// RegisterHTTP register HTTP service to consul
func (c *consulResolver) RegisterHTTP(path, ip string, port int) (string, error) {
	// health check
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("https://%s:%d%s", ip, port, path),
		Method:                         "GET",
		TLSSkipVerify:                  true,
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}

	serviceName := "http" + strings.ReplaceAll(strings.TrimSuffix(path, "/healthz"), "/", ".")
	serviceID := fmt.Sprintf("%s-%s:%d", serviceName, ip, port)
	srv := &api.AgentServiceRegistration{
		ID:      serviceID,                                   // service unique ID
		Name:    serviceName,                                 // service name
		Tags:    []string{utils.GetEnv(), "pudding", "HTTP"}, // service tags
		Address: ip,
		Port:    port,
		Check:   check,
	}

	if err := c.client.Agent().ServiceRegister(srv); err != nil {
		log.Errorf("register http service [%s] failed: %v", serviceID, err)
		return "", err
	}

	log.Infof("register http service [%s] successfully", serviceID)

	return serviceID, nil
}

// Deregister deregister service from consul with serviceID
func (c *consulResolver) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}
