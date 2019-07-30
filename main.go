package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var h2cProxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {
		target, ok := roundrobin(req.Host)
		req.URL.Host = ""
		req.URL.Scheme = "https"
		if ok {
			req.URL.Host = target + ":50051"
		}
	},
	Transport: &http2.Transport{
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	},
}

var proxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {
		target, ok := roundrobin(req.Host)
		req.URL.Host = ""
		req.URL.Scheme = "http"
		if ok {
			req.URL.Host = target
		}
	},
}

func main() {
	var kubeconfig, master string
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()
	klog.InitFlags(nil)
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	server := &http.Server{
		Addr:    ":50051",
		Handler: h2c.NewHandler(http.HandlerFunc(handler), &http2.Server{}),
	}
	go startCtl(config)
	klog.Fatal(server.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		h2cProxy.ServeHTTP(w, r)
	} else {
		proxy.ServeHTTP(w, r)
	}
}
