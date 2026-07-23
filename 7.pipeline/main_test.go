package main

import (
	"slices"
	"testing"
	"time"
)

func TestPipelineSquaresValuesInSourceOrder(t *testing.T) {
	done := make(chan struct{})
	source := sliceToChannel(done, []int{1, 2, 3, 4})
	output := sq(done, source)

	got := collectPipeline(t, output)
	want := []int{1, 4, 9, 16}
	if !slices.Equal(got, want) {
		t.Fatalf("pipeline output = %v; want %v", got, want)
	}

	waitForPipelineClose(t, source)
	close(done)
}

func TestPipelineStopsWhenDownstreamCancels(t *testing.T) {
	done := make(chan struct{})
	source := sliceToChannel(done, []int{1, 2, 3, 4})
	output := sq(done, source)

	close(done)
	waitForPipelineClose(t, output)
	waitForPipelineClose(t, source)
}

func collectPipeline(t *testing.T, input <-chan int) []int {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []int
	for {
		select {
		case value, open := <-input:
			if !open {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("pipeline did not close")
			return nil
		}
	}
}

func waitForPipelineClose(t *testing.T, input <-chan int) {
	t.Helper()
	_ = collectPipeline(t, input)
}
