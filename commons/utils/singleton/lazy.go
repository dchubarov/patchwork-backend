package singleton

import "sync"

// Lazy contains an object that is initialized only once, and provides access to its instance
type Lazy[T any] struct {
	once     sync.Once
	factory  func() T
	instance T
}

// Instance returns the singleton's underlying instance, initialized once with factory
func (s *Lazy[T]) Instance() T {
	s.once.Do(func() {
		if s.factory != nil {
			s.instance = s.factory()
		}
	})
	return s.instance
}

// NewLazy creates new singleton with lazy initialization
func NewLazy[T any](factory func() T) Lazy[T] {
	return Lazy[T]{factory: factory}
}
