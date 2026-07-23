package main

import (
	"slices"
	"sync"
	"testing"
	"time"
)

func TestPartitionedWorkersWriteOrderedResults(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	result := make([]int, len(nums))
	var wg sync.WaitGroup

	for i, num := range nums {
		wg.Go(func() {
			processNum(num, &result[i])
		})
	}
	waitForPartitionedWorkers(t, &wg)

	want := []int{2, 4, 6, 8, 10}
	if !slices.Equal(result, want) {
		t.Fatalf("partitioned results = %v; want %v", result, want)
	}
}

func waitForPartitionedWorkers(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("partitioned workers did not finish")
	}
}
