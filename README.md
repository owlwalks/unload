# Unload

[![GoDoc](https://godoc.org/github.com/owlwalks/unload?status.svg)](https://godoc.org/github.com/owlwalks/unload)
[![Build Status](https://travis-ci.com/owlwalks/unload.svg?branch=master)](https://travis-ci.com/owlwalks/unload)

Minimal gRPC load balancer (h2c - h2 without TLS) that can be used for internal traffic.

_The author has repurposed the life of this project, it was a generic load balancer originally_.

Notes: The heavy lifting is done by [h2c](https://github.com/golang/net/tree/master/http2/h2c), if you want to dive into protocol specifics (negotiation, stream, frames processing or data encoding) please head [here](https://github.com/golang/net/tree/master/http2).