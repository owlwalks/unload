package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var proxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {
		target, ok := roundrobin(req.URL.Host)
		if ok {
			req.URL.Scheme = "https"
			req.URL.Host = target
		}
	},
	Transport: &http2.Transport{
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	},
}

func main() {
	log.SetFlags(log.Llongfile)
	server := &http.Server{
		Addr:    ":50051",
		Handler: h2c.NewHandler(http.HandlerFunc(grpcHandler), &http2.Server{}),
	}
	log.Fatal(server.ListenAndServe())
}

func grpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		proxy.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}
