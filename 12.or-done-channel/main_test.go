package main

import (
	"reflect"
	"testing"
	"time"
)

func TestOrDoneForwardsValuesUntilInputCloses(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int, 3)
	for _, value := range []int{1, 2, 3} {
		input <- value
	}
	close(input)

	got := collectOutput(t, orDone(done, input))

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("orDone returned %v; want %v", got, want)
	}
}

func TestOrDoneStopsWhenCanceledWhileWaitingForInput(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int)
	output := orDone(done, input)

	close(done)
	waitForOutputClose(t, output)
}

func TestOrDoneStopsWhenCanceledWhileOutputIsBlocked(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int)
	_, completed := orDoneWithCompletion(done, input)
	sent := make(chan struct{})

	go func() {
		input <- 42
		close(sent)
	}()

	waitForForwarderSignal(t, sent, "orDone did not receive its input value")
	close(done)
	waitForForwarderSignal(t, completed, "orDone forwarder did not stop after cancellation")
}

func waitForOutputClose[T any](t *testing.T, channel <-chan T) {
	t.Helper()
	_ = collectOutput(t, channel)
}

func collectOutput[T any](t *testing.T, channel <-chan T) []T {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []T
	for {
		select {
		case value, open := <-channel:
			if !open {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("channel did not close after cancellation")
			return nil
		}
	}
}

func waitForForwarderSignal(t *testing.T, signal <-chan struct{}, failure string) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal(failure)
	}
}
