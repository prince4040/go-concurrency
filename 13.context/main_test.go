package main

import (
	"context"
	"testing"
	"time"
)

func TestProducerClosesOutputAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	output := producer(ctx, "value", time.Hour)

	cancel()

	select {
	case _, open := <-output:
		if open {
			t.Fatal("producer sent a value after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("producer did not close output after cancellation")
	}
}

func TestExampleCompletesWithinParentDeadline(t *testing.T) {
	finished := make(chan struct{})

	go func() {
		runExample(40*time.Millisecond, 20*time.Millisecond, 5*time.Millisecond)
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("example did not finish after the parent deadline")
	}
}
