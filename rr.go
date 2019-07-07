package main

import (
	"sync"
	"sync/atomic"
)

var (
	gResolv = struct {
		resolv map[string][]string
		index  map[string]index
		sync.RWMutex
	}{
		resolv: map[string][]string{},
		index:  map[string]index{},
	}
	gNext = []int64{}
)

type index struct {
	idx int
	max int64
}

func roundrobin(host string) (string, bool) {
	var (
		idx     index
		ok      bool
		targets []string
	)
	gResolv.RLock()
	idx, ok = gResolv.index[host]
	targets = gResolv.resolv[host]
	gResolv.RUnlock()
	var next string
	if !ok || idx.max < 0 {
		return next, false
	}
	if idx.max == 0 {
		return targets[0], true
	}
	incr := atomic.AddInt64(&gNext[idx.idx], 1)
	if incr > idx.max {
		incr = 0
		atomic.StoreInt64(&gNext[idx.idx], -1)
	} else {
		atomic.StoreInt64(&gNext[idx.idx], incr)
	}
	return targets[incr], true
}
