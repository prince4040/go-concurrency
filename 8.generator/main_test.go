package main

import (
	"reflect"
	"testing"
	"time"
)

func collect[T any](t *testing.T, stream <-chan T) []T {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []T
	for {
		select {
		case value, open := <-stream:
			if !open {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("stream did not close")
			return nil
		}
	}
}

func TestTakeLimitsOutput(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int, 4)
	for _, value := range []int{1, 2, 3, 4} {
		input <- value
	}
	close(input)

	got := collect(t, take(done, input, 2))
	want := []int{1, 2}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("take returned %v; want %v", got, want)
	}
}

func TestTakeStopsWhenInputCloses(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int, 2)
	input <- 10
	input <- 20
	close(input)

	got := collect(t, take(done, input, 5))
	want := []int{10, 20}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("take returned %v; want %v", got, want)
	}
}

func TestTakeStopsWhenCanceledWhileReceiving(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int)
	output := take(done, input, 1)

	close(done)
	waitForClose(t, output)
}

func TestTakeStopsWhenCanceledWhileSending(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int)
	_, completed := takeWithCompletion(done, input, 1)
	sent := make(chan struct{})

	go func() {
		input <- 42
		close(sent)
	}()

	waitForSignal(t, sent, "take did not receive its input value")
	close(done)
	waitForSignal(t, completed, "take did not stop after cancellation")
}

func waitForClose[T any](t *testing.T, channel <-chan T) {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case _, open := <-channel:
			if !open {
				return
			}
		case <-timer.C:
			t.Fatal("channel did not close after cancellation")
		}
	}
}

func waitForSignal(t *testing.T, signal <-chan struct{}, failure string) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal(failure)
	}
}
