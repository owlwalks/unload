package main

func roundrobin(host string) (string, bool) {
	var next string
	gResolv.Lock()
	defer gResolv.Unlock()
	r, ok := gResolv.resolv[host]
	if !ok {
		return "", false
	}
	next = r.value
	gResolv.resolv[host] = r.next()
	return next, true
}
