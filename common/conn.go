package common

import (
	"container/list"
	"net"
)

type Connection struct {
	fromConn     *net.Conn
	toConn       *net.Conn
	sendHandlers *list.List
	recvHandlers *list.List
}
