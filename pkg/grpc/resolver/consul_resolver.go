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
		GRPC:     fmt.Sprintf("%s:%d", ip, port),
		Timeout:  "10s",
		Interval: "10s",
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
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
		HTTP:     fmt.Sprintf("http://%s:%d%s", ip, port, path),
		Method:   "GET",
		Timeout:  "10s",
		Interval: "10s",
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
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

func (c *consulResolver) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}
