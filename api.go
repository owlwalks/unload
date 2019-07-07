package main

import (
	"encoding/json"
	"net/http"
)

type (
	conf struct {
		Host    string   `json:"host"`
		Targets []string `json:"targets"`
	}
	resErr struct {
		Err string `json:"err"`
	}
)

func setConf(w http.ResponseWriter, r *http.Request) {
	var cfs []conf
	err := json.NewDecoder(r.Body).Decode(&cfs)
	if err != nil {
		_ = json.NewEncoder(w).Encode(resErr{err.Error()})
	}
	gResolv.Lock()
	for _, c := range cfs {
		if _, ok := gResolv.resolv[c.Host]; !ok {
			gResolv.index[c.Host] = index{
				idx: len(gResolv.resolv) - 1,
				max: int64(len(c.Targets)) - 1,
			}
			gNext = append(gNext, -1)
		}
		gResolv.resolv[c.Host] = c.Targets
	}
	gResolv.Unlock()
	_ = json.NewEncoder(w).Encode(struct{}{})
}

func info(w http.ResponseWriter, r *http.Request) {
	var resolv map[string][]string
	gResolv.RLock()
	resolv = gResolv.resolv
	gResolv.RUnlock()
	_ = json.NewEncoder(w).Encode(resolv)
}
