package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"time"

	"github.com/owlwalks/unload"
)

func main() {
	// routing rule, extract service name from uri and host, this is arbitrary
	matcher := func(uri, host []byte) string {
		b := bytes.NewBuffer(uri)
		b.ReadBytes('/')
		service, _ := b.ReadBytes('/')

		if len(service) > 0 {
			if service[len(service)-1] == '/' {
				// "/user/actionX" => "user"
				// "/user/"        => "user"
				// "/user"         => "user"
				service = service[:len(service)-1]
			}
			return string(service)
		}

		return ""
	}

	sc := unload.NewScheduler(true, 10*time.Second, nil)
	p := unload.NewProxy(matcher)
	p.Sch = sc

	// local self-signed: "openssl genrsa -out key 2048"
	key := []byte{}
	// local self-signed: "openssl req -new -x509 -sha256 -key key -out cert -days 365"
	cert := []byte{}

	crt, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{crt}}
	p.ListenTLS(8090, cfg)
}
