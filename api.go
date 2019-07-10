package main

import (
	"sync"
)

var (
	gResolv = struct {
		resolv map[string][]string
		dst    map[string]struct{}
		index  map[string]int
		sync.RWMutex
	}{
		resolv: map[string][]string{},
		dst:    map[string]struct{}{},
		index:  map[string]int{},
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
	if _, ok := gResolv.dst[c.dst]; !ok {
		return
	}
	delete(gResolv.dst, c.dst)
	for i, target := range gResolv.resolv[c.host] {
		if target == c.target {
			gResolv.resolv[c.host] = append(gResolv.resolv[c.host][:i], gResolv.resolv[c.host][i+1:]...)
			gResolv.index[c.host] = -1
			break
		}
	}
}

func addDst(c conf) {
	gResolv.Lock()
	defer gResolv.Unlock()
	if _, ok := gResolv.dst[c.dst]; !ok {
		gResolv.dst[c.dst] = struct{}{}
	}
	if _, ok := gResolv.resolv[c.host]; !ok {
		gResolv.resolv[c.host] = []string{c.target}
	} else {
		gResolv.resolv[c.host] = append(gResolv.resolv[c.host], c.target)
	}
	gResolv.index[c.host] = -1
}
