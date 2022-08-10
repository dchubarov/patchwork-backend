package singleton

import "testing"

type testData struct {
	counter int
}

func (t *testData) inc() *testData {
	t.counter++
	return t
}

func TestSingleton_LazyInit(t *testing.T) {
	s := Lazy(func() *testData {
		return new(testData).inc()
	})

	if s.instance != nil {
		t.Error("Lazy instance initialised before Instance() was invoked")
	}

	testSingleness(t, s)
}

func TestSingleton_EagerInit(t *testing.T) {
	s := Eager(func() *testData {
		return new(testData).inc()
	})

	if s.instance == nil {
		t.Error("Eager singleton not initialised after creation")
	}

	testSingleness(t, s)
}

func testSingleness(t *testing.T, s S[*testData]) {
	s0 := s.Instance()

	for i := 1; i < 1000; i++ {
		s1 := s.Instance()

		if s0 != s1 {
			t.Error("Instances are not equal")
		}

		if s0.counter != 1 || s0.counter != s1.counter {
			t.Error("counters do not match")
		}
	}
}
