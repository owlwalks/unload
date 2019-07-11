package main

import (
	"sync"
)

var (
	gResolv = struct {
		resolv map[string]*ring
		sync.RWMutex
	}{
		resolv: map[string]*ring{},
	}
)

type (
	conf struct {
		host   string
		target string
		dst    string
	}
)

func newConf(host string, target string) conf {
	return conf{
		host:   host,
		target: target,
		dst:    host + target,
	}
}

func rmDst(c conf) {
	gResolv.Lock()
	defer gResolv.Unlock()
	r, ok := gResolv.resolv[c.dst]
	if !ok {
		return
	}
	if r.len() == 1 && r.value == c.target {
		delete(gResolv.resolv, c.dst)
		return
	}
	for p := r.next(); p != r; p = p.nxt {
		if p.value == c.target {
			p.unlink(1)
			break
		}
	}
}

func addDst(c conf) {
	gResolv.Lock()
	defer gResolv.Unlock()
	r, ok := gResolv.resolv[c.dst]
	n := newRing(1)
	n.value = c.target
	if !ok {
		gResolv.resolv[c.host] = n
	} else {
		var skip bool
		// dedup
		for p := r.next(); p != r; p = p.nxt {
			if p.value == c.target {
				skip = true
				break
			}
		}
		if !skip {
			r.link(n)
		}
	}
}
