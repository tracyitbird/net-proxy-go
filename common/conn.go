package common

import (
	"net"
	"container/list"
)

type Connection struct {
	fromConn *net.Conn
	toConn	 *net.Conn
	sendHandlers *list.List
	recvHandlers *list.List
}
