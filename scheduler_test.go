package unload

import (
	"net"
	"testing"
)

func Test_scheduler_requeue(t *testing.T) {
	s := NewScheduler(false, 0, nil)

	name := "_user._tcp."

	s.services[name] = []net.SRV{
		{Target: "user1.staging.local.", Port: 80, Priority: 0, Weight: 30},
		{Target: "user2.staging.local.", Port: 80, Priority: 0, Weight: 20},
		{Target: "user3.staging.local.", Port: 80, Priority: 0, Weight: 40},
	}

	schedule := [...]string{
		"1", "3", "3", "1", "3", "2", "1", "3", "2",
		"1", "3", "3", "1", "3", "2", "1", "3", "2",
	}

	for i := 0; i < 18; i++ {
		next := s.NextBackend(name)
		if string(next.Target[4]) != schedule[i] {
			t.Fail()
		}
	}

	s.services[name] = []net.SRV{
		{Target: "user1.staging.local.", Port: 80, Priority: 0, Weight: 30},
		{Target: "user2.staging.local.", Port: 80, Priority: 0, Weight: 20},
		{Target: "user3.staging.local.", Port: 80, Priority: 10, Weight: 40},
	}

	schedule = [...]string{
		"1", "2", "1", "2", "1", "3", "3", "3", "3",
		"1", "2", "1", "2", "1", "3", "3", "3", "3",
	}

	for i := 0; i < 18; i++ {
		next := s.NextBackend(name)
		if string(next.Target[4]) != schedule[i] {
			t.Fail()
		}
	}
}
