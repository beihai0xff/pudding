// Package utils provides some common utils.
package utils

import (
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/beihai0xff/pudding/pkg/log"
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
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalf("failed to get outbound ip: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}

// GetRand get a random number in [min, max).
func GetRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

// GetHealthEndpointPath get health check http endpoint path.
func GetHealthEndpointPath(prefix string) string {
	return prefix + "/healthz"
}

// GetSwaggerEndpointPath get Swagger ui http endpoint path.
func GetSwaggerEndpointPath(prefix string) string {
	return prefix + "/swagger"
}
