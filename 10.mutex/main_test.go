package main

import (
	"slices"
	"sync"
	"testing"
	"time"
)

func TestProcessNumSafelyAppendsEveryResult(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	var (
		lock   sync.Mutex
		wg     sync.WaitGroup
		result []int
	)

	for _, num := range nums {
		wg.Go(func() {
			processNum(&lock, num, &result)
		})
	}
	waitForWorkers(t, &wg)
	slices.Sort(result)

	want := []int{2, 4, 6, 8, 10}
	if !slices.Equal(result, want) {
		t.Fatalf("processed results = %v; want %v", result, want)
	}
}

func waitForWorkers(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()

	finished := make(chan struct{})
	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("workers did not finish")
	}
}
