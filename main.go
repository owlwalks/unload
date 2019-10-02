package main

import (
	"crypto/tls"
	"flag"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/google/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	resolves       mstring
	resolvesMap    = make(map[string]string)
	targetGroupArn = flag.String("tg-arn", "", "ec2 target group arn")
)

var h2cProxy = &httputil.ReverseProxy{
	Director: func(req *http.Request) {
		target, ok := resolve(req.Host)
		req.URL.Host = ""
		req.URL.Scheme = "https"
		if ok {
			req.URL.Host = target
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
		target, ok := resolve(req.Host)
		req.URL.Host = ""
		req.URL.Scheme = "http"
		if ok {
			req.URL.Host = target
		}
	},
}

func main() {
	logger.Init("", false, false, os.Stderr)
	rand.Seed(time.Now().UnixNano())
	flag.Var(&resolves, "resolve", "Resolve x.staging.service to x.default.svc.cluster.local: -resolve=staging.service,default.svc.cluster.local")
	flag.Parse()
	for _, r := range resolves {
		splits := strings.Split(r, ",")
		if len(splits) > 1 {
			resolvesMap[strings.TrimSpace(splits[0])] = strings.TrimSpace(splits[1])
		}
	}
	go startCtl()
	server := &http.Server{
		Addr: ":50051",
		Handler: h2c.NewHandler(http.HandlerFunc(handler), &http2.Server{
			IdleTimeout: 60 * time.Second,
		}),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Fatal(server.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		h2cProxy.ServeHTTP(w, r)
	} else {
		proxy.ServeHTTP(w, r)
	}
}

func resolve(hostport string) (target string, ok bool) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		// default fallback
		host = hostport
		port = "50051"
	}
	for k, v := range resolvesMap {
		if strings.HasSuffix(host, k) {
			sub := strings.TrimSuffix(host, k)
			host = sub + v
			break
		}
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		logger.Errorln(err)
		return "", false
	}
	if len(ips) == 0 {
		logger.Warningf("%s not found", host)
		return "", false
	}
	var index int
	if len(ips) > 1 {
		index = random(0, len(ips)-1)
	}
	return ips[index].String() + ":" + port, true
}

func random(min, max int) int {
	return min + rand.Intn(max-min)
}

type mstring []string

func (m *mstring) String() string {
	return ""
}

func (m *mstring) Set(val string) error {
	*m = append(*m, val)
	return nil
}
