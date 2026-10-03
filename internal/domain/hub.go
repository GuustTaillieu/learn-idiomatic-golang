package domain

type Hub[T any] interface {
	Subscribe() chan T
	Unsubscribe(ch chan T)
	Broadcast(event T)
}

type Broadcaster[T any] interface {
	Broadcast(event T)
}
