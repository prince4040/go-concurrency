package main

import (
	"slices"
	"testing"
	"time"
)

func TestRunJobsWaitsForEveryJob(t *testing.T) {
	ids := []int{1, 2, 3}
	completed := make(chan int, len(ids))

	finished := make(chan struct{})
	go func() {
		runJobs(ids, func(id int) {
			completed <- id
		})
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("runJobs did not return after every callback completed")
	}
	close(completed)

	var got []int
	for id := range completed {
		got = append(got, id)
	}
	slices.Sort(got)

	if !slices.Equal(got, ids) {
		t.Fatalf("runJobs completed %v; want each of %v", got, ids)
	}
}
