package singleton

import "sync"

// S contains an object that is initialized only once, and provides access to its instance
type S[T any] struct {
	once     sync.Once
	factory  func() T
	instance T
}

// Instance returns the singleton's underlying instance, initialized once with factory
func (s *S[T]) Instance() T {
	s.once.Do(func() {
		if s.factory != nil {
			s.instance = s.factory()
		}
	})
	return s.instance
}

// Lazy creates new singleton with lazy initialization
func Lazy[T any](factory func() T) S[T] {
	return S[T]{factory: factory}
}

// Eager creates new singleton with eager initialization
func Eager[T any](factory func() T) S[T] {
	s := Lazy(factory)
	s.Instance()
	return s
}
