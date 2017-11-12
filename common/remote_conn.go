package common

import (
	"net"
)

func GetRemoteConn(addr string, port string) (net.Conn, error) {
	return NewRemoteConn(addr, port)
}

func NewRemoteConn(addr string, port string) (net.Conn, error) {
	addrAndPort := addr + ":" + port
	return net.Dial("tcp", addrAndPort)
}
