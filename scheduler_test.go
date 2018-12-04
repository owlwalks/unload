package unload

import (
	"net"
	"testing"
)

func Test_scheduler_requeue(t *testing.T) {
	s := NewScheduler(false, 0)

	name := "_user._tcp."

	s.services[name] = []net.SRV{
		{"user1.staging.local.", 80, 0, 30},
		{"user2.staging.local.", 80, 0, 20},
		{"user3.staging.local.", 80, 0, 40},
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
		{"user1.staging.local.", 80, 0, 30},
		{"user2.staging.local.", 80, 0, 20},
		{"user3.staging.local.", 80, 10, 40},
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
