package hook

import (
	"errors"
	"sync"
)

var StopPropagation = errors.New("Event hook propagation stopped")

// Handler defines a hook handler function.
type Handler[T any] func(e T) error

// Hook defines a concurrent safe structure for handling event hooks
// (aka. callbacks propagation).
type Hook[T any] struct {
	mux      sync.RWMutex
	handlers []Handler[T]
}

func (h *Hook[T]) PreAdd(fn Handler[T]) {
	h.mux.Lock()
	defer h.mux.Unlock()

	// minimize allocations by shifting the slice
	h.handlers = append(h.handlers, nil)
	copy(h.handlers[1:], h.handlers)
	h.handlers[0] = fn
}

func (h *Hook[T]) Add(fn Handler[T]) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = append(h.handlers, fn)
}

func (h *Hook[T]) Reset() {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.handlers = nil
}

func (h *Hook[T]) Trigger(data T, oneOffHandlers ...Handler[T]) error {
	h.mux.RLock()
	handlers := make([]Handler[T], 0, len(h.handlers)+len(oneOffHandlers))
	handlers = append(handlers, h.handlers...)
	handlers = append(handlers, oneOffHandlers...)
	// unlock is not deferred to avoid deadlocks when Trigger is called recursive by the handlers
	h.mux.RUnlock()

	for _, fn := range handlers {
		err := fn(data)
		if err == nil {
			continue
		}

		if errors.Is(err, StopPropagation) {
			return nil
		}

		return err
	}

	return nil
}
