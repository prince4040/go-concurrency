package main

import (
	"slices"
	"testing"
	"time"
)

func TestFibonacciProducesPrefixAndStopsAfterCancellation(t *testing.T) {
	done := make(chan struct{})
	numbers := fibonacci(done)

	got := make([]int, 0, 10)
	for range 10 {
		select {
		case value, open := <-numbers:
			if !open {
				t.Fatal("fibonacci output closed before ten values")
			}
			got = append(got, value)
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for a Fibonacci value")
		}
	}

	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	if !slices.Equal(got, want) {
		t.Fatalf("fibonacci prefix = %v; want %v", got, want)
	}

	close(done)
	waitForIntChannelClose(t, numbers)
}

func waitForIntChannelClose(t *testing.T, input <-chan int) {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case _, open := <-input:
			if !open {
				return
			}
		case <-timer.C:
			t.Fatal("channel did not close after cancellation")
		}
	}
}
