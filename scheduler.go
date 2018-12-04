package unload

import (
	"container/heap"
	"net"
	"sort"
	"sync"
	"time"
)

// Lookup is custom DNS lookup function
// Use for static routing
type Lookup func(service string) []net.SRV

// Scheduler schedule backends in weighted round-robin
type Scheduler struct {
	sync.Mutex
	// name in _service._proto.name.
	name             string
	backends         map[string]*queue
	services         map[string][]net.SRV
	Relookup         bool
	RelookupInterval time.Duration
	CustomLookup     Lookup
}

type queue []net.SRV

// NewScheduler makes a new scheduler
// It also starts SRV updating periodically if relookup is set to true
func NewScheduler(relookup bool, interval time.Duration, custom Lookup) *Scheduler {
	s := Scheduler{
		backends:         make(map[string]*queue),
		services:         make(map[string][]net.SRV),
		Relookup:         relookup,
		RelookupInterval: interval,
		CustomLookup:     custom,
	}

	if relookup {
		go s.relookupEvery(interval)
	}

	return &s
}

// NextBackend returns next backend in queue
func (s *Scheduler) NextBackend(service string) net.SRV {
	q, ok := s.getQueue(service)
	if !ok {
		s.lookup(service)
	}

	if q == nil || q.Len() == 0 {
		s.requeue(service)
	}

	return s.pop(service)
}

func (s *Scheduler) pop(service string) net.SRV {
	s.Lock()
	defer s.Unlock()
	q, ok := s.backends[service]
	if ok && q.Len() > 0 {
		return heap.Pop(q).(net.SRV)
	}

	return net.SRV{}
}

func (s *Scheduler) getQueue(service string) (q *queue, ok bool) {
	s.Lock()
	q, ok = s.backends[service]
	s.Unlock()
	return
}

func (s *Scheduler) getSRVs(service string) (backends []net.SRV, ok bool) {
	s.Lock()
	backends, ok = s.services[service]
	s.Unlock()
	return
}

func (s *Scheduler) requeue(service string) {
	records, _ := s.getSRVs(service)
	nRecords := len(records)
	if nRecords == 0 {
		return
	}

	total := uint(0)
	for _, val := range records {
		total += uint(val.Weight)
	}

	unordered := make([]int, nRecords)
	for i, val := range records {
		pct := 1.0
		if total != 0 {
			pct = float64(val.Weight) / float64(total) * 10
		}
		unordered[i] = int(pct)
	}

	ordered := append(unordered[:0:0], unordered...)
	sort.Ints(ordered)

	q := queue{}
	max := ordered[nRecords-1]
	for rep := 1; rep <= max; rep++ {
		for index := 0; index < nRecords; index++ {
			if unordered[index]-rep >= 0 {
				q = append(q, records[index])
			}
		}
	}

	ptr := &q
	heap.Init(ptr)
	s.Lock()
	s.backends[service] = ptr
	s.Unlock()
}

func (s *Scheduler) lookup(service string) error {
	var records []net.SRV
	if s.CustomLookup != nil {
		records = s.CustomLookup(service)
	} else {
		_, addrs, err := net.LookupSRV(service, "tcp", s.name)
		if err != nil {
			return err
		}
		// use values here
		records = make([]net.SRV, len(addrs))
		for i := 0; i < len(addrs); i++ {
			records[i] = *addrs[i]
		}
	}
	s.Lock()
	s.services[service] = records
	s.Unlock()
	return nil
}

func (s *Scheduler) relookupEvery(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.Lock()
			services := make([]string, len(s.services))
			i := 0
			for k := range s.services {
				services[i] = k
				i++
			}
			s.Unlock()
			for _, service := range services {
				go s.lookup(service)
			}
		}
	}
}

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].Priority < q[j].Priority }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(net.SRV))
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}
