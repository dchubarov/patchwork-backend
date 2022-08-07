package shutdown

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestShutdown_AllHandlersLifoOrder(t *testing.T) {
	var order []int

	s := &Shutdown{}
	s.Register("first", time.Second, func(ctx context.Context) error {
		order = append(order, 1)
		return nil
	})
	s.Register("second", time.Second, func(ctx context.Context) error {
		order = append(order, 2)
		return nil
	})
	s.Register("third", time.Second, func(ctx context.Context) error {
		order = append(order, 3)
		return nil
	})

	s.ShutdownAll()

	for _, h := range s.hooks {
		if !h.done {
			t.Errorf("Handler has not been executed: %s", h.tag)
		}
	}

	if !reflect.DeepEqual(order, []int{3, 2, 1}) {
		t.Errorf("Shutdown handlers executed in wrong order: %v", order)
	}
}
