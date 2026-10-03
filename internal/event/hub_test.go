package event_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/event"
)

func TestHub_SubscribersGetBroadcastedEvents(t *testing.T) {
	// Arrange
	hub := event.NewHub[string]()
	subscriber1 := hub.Subscribe()
	subscriber2 := hub.Subscribe()

	// Act
	hub.Broadcast("event1")
	hub.Broadcast("event2")

	// Assert
	for _, subscriber := range []chan string{subscriber1, subscriber2} {
		select {
		case event := <-subscriber:
			if event != "event1" && event != "event2" {
				t.Errorf("expected 'event1' or 'event2', but got %s", event)
			}
		default:
			t.Errorf("expected to receive an event, but channel was empty")
		}
	}
}

func TestHub_UnsubscribeStopsReceivingEvents(t *testing.T) {
	// Arrange
	hub := event.NewHub[string]()
	subscriber := hub.Subscribe()

	// Act
	hub.Unsubscribe(subscriber)
	hub.Broadcast("event1")

	// Assert: channel should be closed and no events should be received
	select {
	case _, ok := <-subscriber:
		if ok {
			t.Errorf("expected channel to be closed, but it was open")
		}
	default:
		t.Errorf("expected channel to be closed, but it was not")
	}
}

func TestHub_BroadcastDoesNotBlock(t *testing.T) {
	// Arrange
	hub := event.NewHub[string]()
	_ = hub.Subscribe()

	// Fill the subscriber channel to its capacity
	for i := 0; i < 16; i++ {
		hub.Broadcast(fmt.Sprintf("event%d", i))
	}

	// Act: Broadcast should not block even if the subscriber channel is full
	done := make(chan struct{})
	go func() {
		hub.Broadcast("eventHE")
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Broadcast did not block as expected")
	}
}
