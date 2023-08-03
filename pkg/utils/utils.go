// Package utils provides some common utils.
package utils

import (
	"math/rand"
	"net"
	"os"
)

const (
	puddingEnv = "PUDDING_ENV"
	defaultEnv = "dev"
)

// GetEnv get env environment variables.
func GetEnv() string {
	env := os.Getenv(puddingEnv)
	if env == "" {
		return defaultEnv
	}

	return env
}

// GetOutBoundIP get preferred outbound ip of this machine.
func GetOutBoundIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String()
		}
	}

	return ""
}

// GetRand get a random number in [min, max).
func GetRand(min, max int) int {
	return rand.Intn(max-min) + min
}

// GetHealthEndpointPath get health check http endpoint path.
func GetHealthEndpointPath(prefix string) string {
	return "/healthz"
}

// GetSwaggerEndpointPath get Swagger ui http endpoint path.
func GetSwaggerEndpointPath(prefix string) string {
	return prefix + "/swagger"
}
