package optional

import "reflect"

type O[T any] struct {
	present bool
	value   T
}

// Present returns true if optional contains a value, otherwise returns false
func (o *O[T]) Present() bool {
	return o.present
}

// Get returns a value of an optional instance, and a flag indicating whether value is present
func (o *O[T]) Get() (T, bool) {
	return o.value, o.present
}

// Else returns a value of optional, or default value passed by argument
func (o *O[T]) Else(v T) T {
	if o.present {
		return o.value
	} else {
		return v
	}
}

// Empty returns an empty optional
func Empty[T any]() *O[T] {
	return &O[T]{present: false}
}

// Of returns an optional containing a specified value, or an empty optional of value is zero
func Of[T any](v T) *O[T] {
	if reflect.ValueOf(v).IsZero() {
		return &O[T]{present: false}
	} else {
		return &O[T]{true, v}
	}
}

// Map invokes mapper function if optional contains a value, returns optional containing mapping result
func Map[T any, U any](v *O[T], mapper func(v T) U) *O[U] {
	if v.present {
		return Of(mapper(v.value))
	} else {
		return Empty[U]()
	}
}

// MapValue is a shortcut for Map(Of(v), mapper)
func MapValue[T any, U any](v T, mapper func(v T) U) *O[U] {
	return Map(Of(v), mapper)
}
