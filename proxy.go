package unload

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	errReadHeaderTimeout = errors.New("reading header timeout")
)

type Proxy struct {
	sync.Mutex
	sc    Scheduler
	conns map[*tcpConn]struct{}
}

func NewProxy() Proxy {
	return Proxy{
		sc:    NewScheduler(),
		conns: make(map[*tcpConn]struct{}),
	}
}

func (p *Proxy) Listen(port int) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	delay := 5 * time.Millisecond
	for {
		conn, e := l.AcceptTCP()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if max := 1 * time.Second; delay > max {
					delay = max
				}
				time.Sleep(delay)
				delay *= 2
				continue
			}
			return
		}
		delay = 5 * time.Millisecond
		src := newConn(conn)
		go p.proxy(src)
	}
}

func (p *Proxy) proxy(src *tcpConn) {
	src.SetKeepAlive(true)
	src.SetKeepAlivePeriod(3 * time.Minute)
	br := newBufioReader(src)
	defer putBufioReader(br)
	var dst *tcpConn
	for {
		header, err := readHeader(src, br)
		if err != nil {
			return
		}
		dst = p.open(&net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: 8090})
		if dst != nil {
			derr := make(chan error)
			uerr := make(chan error)
			dst.Write(header)
			go cp(dst, br, derr)
			go cp(src, dst, uerr)
			for i := 0; i < 2; i++ {
				select {
				case <-derr:
					// downstream is closed, stop reading from upstream
					dst.SetReadDeadline(time.Now())
				case err = <-uerr:
					// upstream is closed, force closing downstream
					if _, ok := err.(*net.OpError); ok {
						p.close(dst)
					}
					p.close(src)
				}
			}
			close(derr)
			close(uerr)
		} else {
			p.close(src)
		}
	}
}

func readHeader(src *tcpConn, br *bufio.Reader) ([]byte, error) {
	tp := newTextprotoReader(br)
	defer putTextprotoReader(tp)

	l1, e := tp.ReadLineBytes()
	if e != nil {
		return nil, e
	}

	b := bytes.NewBuffer(l1)
	b.ReadBytes(' ')
	// first line, between first and second space
	uri, _ := b.ReadBytes(' ')
	if len(uri) > 0 {
		// rm ' ' including from last read
		src.uri = uri[:len(uri)-1]
	}

	l2, e := tp.ReadLineBytes()
	if e != nil {
		return nil, e
	}

	b = bytes.NewBuffer(l2)
	b.ReadBytes(' ')
	src.host, _ = b.ReadBytes('\n')

	l1 = append(l1, byte('\r'), byte('\n'))
	l2 = append(l2, byte('\r'), byte('\n'))

	return append(l1, l2...), nil
}

func (p *Proxy) close(conn *tcpConn) {
	conn.Close()
	p.Lock()
	defer p.Unlock()
	delete(p.conns, conn)
}

func (p *Proxy) open(addr *net.TCPAddr) *tcpConn {
	// need to reuse conn, keep track of its state
	conn, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		return nil
	}
	c := newConn(conn)
	c.SetKeepAlive(true)
	c.SetKeepAlivePeriod(3 * time.Minute)
	return c
}

func cp(dst io.Writer, src io.Reader, result chan error) {
	_, err := io.Copy(dst, src)
	result <- err
}
