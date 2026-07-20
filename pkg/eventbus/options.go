package eventbus

type EventBusOption[T any] func(*EventBus[T])

func WithDropHandler[T any](fn func(T)) EventBusOption[T] {
	return func(b *EventBus[T]) {
		b.onDrop = fn
	}
}
