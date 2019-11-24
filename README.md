# UNMAINTAINED, PLEASE SEE OTHER ALTERNATIVES BELOW

* [ingress-nginx](https://github.com/kubernetes/ingress-nginx)

## Unload

[![GoDoc](https://godoc.org/github.com/owlwalks/unload?status.svg)](https://godoc.org/github.com/owlwalks/unload)
[![Build Status](https://travis-ci.com/owlwalks/unload.svg?branch=master)](https://travis-ci.com/owlwalks/unload)

Minimal gRPC load balancer (h2c - h2 without TLS) that can be used for internal traffic.

_The author has repurposed the life of this project, it was originally a generic load balancer_.

## How-to

* [Deploy to Kubernetes](examples)
* [Develop](docs/dev)

## Notes
For protocol specifics (negotiation, stream, frames processing or data encoding) please head [here](https://github.com/golang/net/tree/master/http2/h2c) and [here](https://github.com/golang/net/tree/master/http2).
