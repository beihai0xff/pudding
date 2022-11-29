package utils

import (
	"net"
	"strings"

	"github.com/beihai0xff/pudding/pkg/log"
)

func GetOutBoundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalf("failed to get outbound ip: %w", err)
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.String(), ":")[0]
}
