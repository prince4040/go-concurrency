package main

import (
	"slices"
	"testing"
	"time"
)

func TestRunNumbersWaitsForEveryNumber(t *testing.T) {
	nums := []string{"1", "2", "3"}
	printed := make(chan string, len(nums))

	finished := make(chan struct{})
	go func() {
		runNumbers(nums, func(num string) {
			printed <- num
		})
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("runNumbers did not return after every callback completed")
	}
	close(printed)

	var got []string
	for num := range printed {
		got = append(got, num)
	}
	slices.Sort(got)

	if !slices.Equal(got, nums) {
		t.Fatalf("runNumbers printed %v; want each of %v", got, nums)
	}
}
