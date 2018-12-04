package main

import (
	"fmt"
	"net"

	"github.com/owlwalks/unload"
)

func main() {
	intro := `Spin up the backends:
docker run -d -p 32768:80 containous/whoami
docker run -d -p 32769:80 containous/whoami

Add 2 records to /etc/hosts:
127.0.0.1 whoami1.local
127.0.0.1 whoami2.local

Run wrk:
wrk -t20 -c1000 -d60s -H "Host: test.traefik" --latency  http://127.0.0.1:8090/bench`

	fmt.Println(intro)

	customLookup := func(service string) []net.SRV {
		return []net.SRV{
			{Target: "whoami1.local", Port: 32768, Weight: 50},
			{Target: "whoami2.local", Port: 32769, Weight: 50},
		}
	}

	matcher := func(uri, host []byte) string {
		return "test.traefik"
	}

	sc := unload.NewScheduler(false, 0, customLookup)
	p := unload.NewProxy(matcher)
	p.Sch = sc

	p.Listen(8090)
}
