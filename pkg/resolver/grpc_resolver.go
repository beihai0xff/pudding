package resolver

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// consulResolver
type consulResolver struct {
	client *api.Client
}

// NewConsulResolver new
func NewConsulResolver(addr string) (Resolver, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	c, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &consulResolver{c}, nil
}

// Register register gRPC service to consul
func (c *consulResolver) Register(serviceName string, ip string, port int) (string, error) {
	// 健康检查
	check := &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", ip, port), // 这里一定是外部可以访问的地址
		Timeout:  "10s",
		Interval: "10s",
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为1分钟，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
		DeregisterCriticalServiceAfter: "1m",
	}

	serviceID := fmt.Sprintf("%s-%s:%d", serviceName, ip, port)
	srv := &api.AgentServiceRegistration{
		ID:      serviceID,           // service unique ID
		Name:    serviceName,         // service name
		Tags:    []string{"pudding"}, // service tags
		Address: ip,
		Port:    port,
		Check:   check,
	}
	return serviceID, c.client.Agent().ServiceRegister(srv)
}
func (c *consulResolver) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}
