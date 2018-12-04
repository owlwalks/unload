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

func (c *tcpConn) Read(b []byte) (int, error) {
	n, err := c.TCPConn.Read(b)
	c.busy = true
	if err != nil {
		c.busy = false
	}
	return n, err
}

func (c *tcpConn) Write(b []byte) (int, error) {
	n, err := c.TCPConn.Write(b)
	c.busy = true
	if err != nil {
		c.busy = false
	}
	return n, err
}

func newConn(c *net.TCPConn) *tcpConn {
	return &tcpConn{c, nil, nil, false}
}
