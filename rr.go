package main

func roundrobin(host string) (string, bool) {
	var (
		idx, max int
		ok       bool
		next     string
		targets  []string
	)
	gResolv.RLock()
	idx, ok = gResolv.index[host]
	targets = make([]string, len(gResolv.resolv[host]))
	copy(targets, gResolv.resolv[host])
	gResolv.RUnlock()
	max = len(targets)
	if !ok || max == 0 {
		return next, false
	}
	if max == 1 {
		return targets[0], true
	}
	idx++
	if idx >= max {
		idx = -1
	}
	gResolv.Lock()
	gResolv.index[host] = idx
	gResolv.Unlock()
	if idx < 0 {
		idx = 0
	}
	return targets[idx], true
}
