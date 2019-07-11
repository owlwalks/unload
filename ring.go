// https://github.com/golang/go/tree/master/src/container/ring

package main

type ring struct {
	nxt, prv *ring
	value    string
}

func (r *ring) init() *ring {
	r.nxt = r
	r.prv = r
	return r
}

func (r *ring) next() *ring {
	if r.nxt == nil {
		return r.init()
	}
	return r.nxt
}

func (r *ring) prev() *ring {
	if r.nxt == nil {
		return r.init()
	}
	return r.prv
}

func (r *ring) move(n int) *ring {
	if r.nxt == nil {
		return r.init()
	}
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prv
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.nxt
		}
	}
	return r
}

func newRing(n int) *ring {
	if n <= 0 {
		return nil
	}
	r := new(ring)
	p := r
	for i := 1; i < n; i++ {
		p.nxt = &ring{prv: p}
		p = p.nxt
	}
	p.nxt = r
	r.prv = p
	return r
}

func (r *ring) link(s *ring) *ring {
	n := r.next()
	if s != nil {
		p := s.prev()
		// Note: Cannot use multiple assignment because
		// evaluation order of LHS is not specified.
		r.nxt = s
		s.prv = r
		n.prv = p
		p.nxt = n
	}
	return n
}

func (r *ring) unlink(n int) *ring {
	if n <= 0 {
		return nil
	}
	return r.link(r.move(n + 1))
}

func (r *ring) len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.next(); p != r; p = p.nxt {
			n++
		}
	}
	return n
}
