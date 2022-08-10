package singleton

import "testing"

type testData struct {
	counter int
}

func TestLazy(t *testing.T) {
	s := Lazy(func() *testData {
		return &testData{0}
	})

	s0 := s.Instance()
	s1 := s.Instance()
	s1.counter++

	if s0 != s1 {
		t.Error("Instances are not equal")
	}

	if s0.counter != 1 || s0.counter != s1.counter {
		t.Error("counters do not match")
	}
}
