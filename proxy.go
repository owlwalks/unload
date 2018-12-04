package unload

import (
	"bufio"
	"bytes"
	"crypto/tls"
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

type Matcher func(uri, host []byte) string

type Proxy struct {
	sync.Mutex
	sc      *Scheduler
	conns   map[string]map[*tcpConn]struct{}
	matcher Matcher
}

func NewProxy(relookup bool, interval time.Duration, matchFn Matcher) *Proxy {
	return &Proxy{
		sc:      NewScheduler(relookup, interval),
		conns:   make(map[string]map[*tcpConn]struct{}),
		matcher: matchFn,
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

func (p *Proxy) ListenTLS(port int, cert, key []byte) {
	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{crt}}
	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	ln := tls.NewListener(l, cfg)
	delay := 5 * time.Millisecond
	for {
		conn, e := ln.Accept()
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
		header, uri, host, err := readHeader(br)
		if err != nil {
			p.close(src)
			return
		}
		addr, err := p.resolve(uri, host)
		if err != nil {
			p.close(src)
			return
		}
		dst = p.open(addr)
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
					p.close(dst)
					p.close(src)
				}
			}
			close(derr)
			close(uerr)
			return
		}
		p.close(src)
		return
	}
}

func (p *Proxy) close(conn *tcpConn) {
	conn.Close()
	saddr := conn.RemoteAddr().String()
	p.Lock()
	defer p.Unlock()
	if _, ok := p.conns[saddr]; ok {
		delete(p.conns[saddr], conn)
	}
}

func (p *Proxy) get(saddr string) *tcpConn {
	p.Lock()
	defer p.Unlock()
	if pool, ok := p.conns[saddr]; ok {
		for conn := range pool {
			if conn.busy {
				continue
			}
			return conn
		}
	}

	return nil
}

func (p *Proxy) open(addr *net.TCPAddr) *tcpConn {
	saddr := addr.String()
	c := p.get(saddr)
	if c != nil {
		// reset deadline because it was set to break io.Copy (blocked) in last routine
		c.SetReadDeadline(time.Time{})
		return c
	}

	conn, e := net.DialTCP("tcp", nil, addr)
	if e != nil {
		return nil
	}

	c = newConn(conn)
	c.SetKeepAlive(true)
	c.SetKeepAlivePeriod(3 * time.Minute)

	p.Lock()
	defer p.Unlock()
	if _, ok := p.conns[saddr]; !ok {
		p.conns[saddr] = make(map[*tcpConn]struct{})
	}
	p.conns[saddr][c] = struct{}{}

	return c
}

func (p *Proxy) resolve(uri, host []byte) (*net.TCPAddr, error) {
	service := p.matcher(uri, host)
	srv := p.sc.NextBackend(service)
	addr, err := net.ResolveTCPAddr("tcp", srv.Target)
	addr.Port = int(srv.Port)

	return addr, err
}

func readHeader(br *bufio.Reader) ([]byte, []byte, []byte, error) {
	tp := newTextprotoReader(br)
	defer putTextprotoReader(tp)

	l1, e := tp.ReadLineBytes()
	if e != nil {
		return nil, nil, nil, e
	}

	b := bytes.NewBuffer(l1)
	b.ReadBytes(' ')
	// first line, between first and second space
	uri, _ := b.ReadBytes(' ')
	if len(uri) > 0 && uri[len(uri)-1] == ' ' {
		// rm ' ' including from last read
		uri = uri[:len(uri)-1]
	}

	l2, e := tp.ReadLineBytes()
	if e != nil {
		return nil, nil, nil, e
	}

	b = bytes.NewBuffer(l2)
	b.ReadBytes(' ')
	host, _ := b.ReadBytes('\n')

	l1 = append(l1, byte('\r'), byte('\n'))
	l2 = append(l2, byte('\r'), byte('\n'))

	return append(l1, l2...), uri, host, nil
}

func cp(dst io.Writer, src io.Reader, result chan error) {
	_, err := io.Copy(dst, src)
	result <- err
}
