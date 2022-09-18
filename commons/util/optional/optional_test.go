package optional

import (
	"testing"
)

func TestOptional_Empty(t *testing.T) {
	o1 := Empty[uint16]()
	if o1.Present() {
		t.Fail()
	}

	o2 := Of("")
	if o2.Present() {
		t.Fail()
	}

	o3 := Of(0.0)
	if o3.Present() {
		t.Fail()
	}

	if o3.Else(5.) != 5. {
		t.Fail()
	}

	o4 := Of[*int](nil)
	if o4.Present() {
		t.Fail()
	}
}

func TestOptional_Mappers(t *testing.T) {
	o1 := MapValue(5, func(v int) float64 {
		return float64(v)
	})

	if o1.Else(0.) < 0. {
		t.Fail()
	}
}
