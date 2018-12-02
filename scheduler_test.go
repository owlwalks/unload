package unload

import (
	"testing"
)

func Test_scheduler_requeue(t *testing.T) {
	s := NewScheduler()

	name := "_user._tcp.staging.local."

	s.discovered[name] = []backend{
		{"user1.staging.local.", 0, 30},
		{"user2.staging.local.", 0, 20},
		{"user3.staging.local.", 0, 40},
	}

	schedule := [...]string{
		"1", "3", "3", "1", "3", "2", "1", "3", "2",
		"1", "3", "3", "1", "3", "2", "1", "3", "2",
	}

	for i := 0; i < 18; i++ {
		target := s.NextBackend(name)
		if string(target[4]) != schedule[i] {
			t.Fail()
		}
	}

	s.discovered[name] = []backend{
		{"user1.staging.local.", 0, 30},
		{"user2.staging.local.", 0, 20},
		{"user3.staging.local.", 10, 40},
	}

	schedule = [...]string{
		"1", "2", "1", "2", "1", "3", "3", "3", "3",
		"1", "2", "1", "2", "1", "3", "3", "3", "3",
	}

	for i := 0; i < 18; i++ {
		target := s.NextBackend(name)
		if string(target[4]) != schedule[i] {
			t.Fail()
		}
	}
}
