package unload

import (
	"net"
	"time"
)

type tcpConn struct {
	rwc  net.Conn
	busy bool
}

func (c *tcpConn) Read(b []byte) (int, error) {
	n, err := c.rwc.Read(b)
	if c.busy = true; err != nil {
		c.busy = false
	}
	return n, err
}

func (c *tcpConn) Write(b []byte) (int, error) {
	n, err := c.rwc.Write(b)
	if c.busy = true; err != nil {
		c.busy = false
	}
	return n, err
}

func (c *tcpConn) Close() error {
	return c.rwc.Close()
}

func (c *tcpConn) SetReadDeadline(t time.Time) error {
	return c.rwc.SetReadDeadline(t)
}

func (c *tcpConn) SetKeepAlive(keepalive bool) error {
	if conn, ok := c.rwc.(*net.TCPConn); ok {
		return conn.SetKeepAlive(keepalive)
	}
	return nil
}

func (c *tcpConn) SetKeepAlivePeriod(d time.Duration) error {
	if conn, ok := c.rwc.(*net.TCPConn); ok {
		return conn.SetKeepAlivePeriod(d)
	}
	return nil
}

func (c *tcpConn) SetLinger(sec int) error {
	if conn, ok := c.rwc.(*net.TCPConn); ok {
		return conn.SetLinger(sec)
	}
	return nil
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

func newConn(c net.Conn) *tcpConn {
	return &tcpConn{c, false}
}
