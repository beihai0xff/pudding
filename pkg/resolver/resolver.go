package resolver

type Resolver interface {
	Register(serviceName string, ip string, port int) (string, error)
	Deregister(serviceID string) error
}
