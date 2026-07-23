package main

import (
	"slices"
	"testing"
	"time"
)

type observedValue struct {
	kind  string
	value string
}

func TestConsumeStreamsPreservesPerStreamOrder(t *testing.T) {
	letters := make(chan string)
	numbers := make(chan string)
	go sendValues(letters, "a", "b", "c")
	go sendValues(numbers, "1", "2", "3")

	finished := make(chan []observedValue, 1)
	go func() {
		var observed []observedValue
		consumeStreams(letters, numbers, func(kind, value string) {
			observed = append(observed, observedValue{kind: kind, value: value})
		})
		finished <- observed
	}()

	var observed []observedValue
	select {
	case observed = <-finished:
	case <-time.After(time.Second):
		t.Fatal("consumeStreams did not return after both inputs closed")
	}

	var gotLetters, gotNumbers []string
	for _, item := range observed {
		switch item.kind {
		case "letter":
			gotLetters = append(gotLetters, item.value)
		case "number":
			gotNumbers = append(gotNumbers, item.value)
		default:
			t.Fatalf("consumeStreams reported unknown kind %q", item.kind)
		}
	}

	if want := []string{"a", "b", "c"}; !slices.Equal(gotLetters, want) {
		t.Fatalf("letter order = %v; want %v", gotLetters, want)
	}
	if want := []string{"1", "2", "3"}; !slices.Equal(gotNumbers, want) {
		t.Fatalf("number order = %v; want %v", gotNumbers, want)
	}
}
