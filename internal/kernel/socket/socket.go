package socket

import (
	"fmt"
	"net"

	"github.com/kodmain/KitsuneFramework/internal/env"
)

const (
	ADDR_TYPE = "unixpacket"
)

var (
	SOCK_SUPERVISOR = env.PATH_RUN + env.BUILD_APP_NAME
)

func Server(socket string) *server {
	sock := &net.UnixAddr{
		Name: socket,
		Net:  ADDR_TYPE,
	}

	fmt.Println(sock)

	return nil
}

func Client() *client {
	return nil
}
