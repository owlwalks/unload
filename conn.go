package unload

import (
	"net"
)

type tcpConn struct {
	*net.TCPConn
	uri  []byte
	host []byte
	busy bool
}

func newConn(c *net.TCPConn) *tcpConn {
	return &tcpConn{c, nil, nil, false}
}
