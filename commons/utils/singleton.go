package utils

import "sync"

// Singleton contains an object that is initialized only once, and provides access to its instance
type Singleton[T any] struct {
	once     sync.Once
	instance T
}

// Instance returns the singleton's underlying instance, initialized once with factory
func (s *Singleton[T]) Instance(factory func() T) T {
	s.once.Do(func() {
		s.instance = factory()
	})
	return s.instance
}
