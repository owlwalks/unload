package unload

import (
	"container/heap"
	"sort"
	"sync"
)

type backend struct {
	target   string
	priority uint
	weight   uint
}

type queue []backend

type Scheduler struct {
	sync.Mutex
	backends   map[string]*queue
	discovered map[string][]backend
}

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].priority < q[j].priority }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(backend))
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

func (s *Scheduler) NextBackend(name string) string {
	s.Lock()
	defer s.Unlock()
	q := s.backends[name]
	if q == nil || q.Len() == 0 {
		s.requeue(name)
	}
	q = s.backends[name]
	if q != nil && q.Len() > 0 {
		b := heap.Pop(q).(backend)
		return b.target
	}
	return ""
}

func (s *Scheduler) requeue(name string) {
	backends := s.discovered[name]
	nBackends := len(backends)
	if nBackends == 0 {
		return
	}

	total := uint(0)
	for _, val := range backends {
		total += val.weight
	}

	unordered := make([]int, nBackends)
	for i, val := range backends {
		pct := 1.0
		if total != 0 {
			pct = float64(val.weight) / float64(total) * 10
		}
		unordered[i] = int(pct)
	}

	ordered := append(unordered[:0:0], unordered...)
	sort.Ints(ordered)

	q := queue{}
	max := ordered[nBackends-1]
	for rep := 1; rep <= max; rep++ {
		for index := 0; index < nBackends; index++ {
			if unordered[index]-rep >= 0 {
				q = append(q, backends[index])
			}
		}
	}

	ptr := &q
	heap.Init(ptr)
	s.backends[name] = ptr
}

func NewScheduler() Scheduler {
	return Scheduler{backends: make(map[string]*queue), discovered: make(map[string][]backend)}
}
