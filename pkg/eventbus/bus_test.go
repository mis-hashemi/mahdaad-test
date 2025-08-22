package eventbus

import (
	"testing"
	"time"
)

func TestPublishSubscribe(t *testing.T) {
	bus := New()

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	bus.Subscribe("course_created", func(e Event) {
		ch1 <- e.Data.(string)
	})

	bus.Subscribe("course_created", func(e Event) {
		ch2 <- e.Data.(string)
	})

	bus.Publish(Event{Name: "course_created", Data: "course-1"})

	select {
	case msg := <-ch1:
		if msg != "course-1" {
			t.Errorf("expected course-1, got %s", msg)
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("timeout waiting for ch1")
	}

	select {
	case msg := <-ch2:
		if msg != "course-1" {
			t.Errorf("expected course-1, got %s", msg)
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("timeout waiting for ch2")
	}
}

func TestPublishWithoutSubscriber(t *testing.T) {
	bus := New()

	// should not panic even if no subscriber
	bus.Publish(Event{Name: "non_existing_event", Data: "course1"})
}

func TestMultipleEvents(t *testing.T) {
	bus := New()
	results := make([]string, 0, 2)
	done := make(chan struct{})

	bus.Subscribe("event1", func(e Event) {
		results = append(results, e.Data.(string))
		if len(results) == 2 {
			close(done)
		}
	})

	bus.Publish(Event{Name: "event1", Data: "first"})
	bus.Publish(Event{Name: "event1", Data: "second"})

	select {
	case <-done:
		foundFirst := false
		foundSecond := false
		for _, r := range results {
			if r == "first" {
				foundFirst = true
			}
			if r == "second" {
				foundSecond = true
			}
		}
		if !foundFirst || !foundSecond {
			t.Errorf("expected both first and second, got %v", results)
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("timeout waiting for events")
	}
}
