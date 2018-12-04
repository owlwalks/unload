# Unload

[![GoDoc](https://godoc.org/github.com/owlwalks/unload?status.svg)](https://godoc.org/github.com/owlwalks/unload)
[![Build Status](https://travis-ci.com/owlwalks/unload.svg?branch=master)](https://travis-ci.com/owlwalks/unload)

*Unload* is an application (level 7) load balancer written in Go. Aiming for simplicity and throughput.

### Example:
```golang
sc := unload.NewScheduler(true, 10*time.Second, nil)
p := unload.NewProxy(nil)
p.ListenTLS(443, cfg)
```
Complete example is [here](https://github.com/owlwalks/unload/blob/master/unload/main.go).

### Features:
  - [x] Service discovery (SRV records)
  - [x] Dynamic routing
  - [x] Static routing
  - [x] Weighted round-robin
  - [x] TCP pooling on backend side
  - [x] TLS

### Todo:
- [ ] HTTP/2
- [ ] Healthcheck

### Benchmark:
* Traefik - [here](https://github.com/owlwalks/unload/tree/master/bench)
* gobetween - Todo