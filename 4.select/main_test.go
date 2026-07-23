package main

import (
	"slices"
	"testing"
	"time"
)

func TestReceiveBothReturnsEveryValue(t *testing.T) {
	finished := make(chan []string, 1)
	go func() {
		finished <- receiveBoth(produce("1"), produce("2"))
	}()

	var got []string
	select {
	case got = <-finished:
	case <-time.After(time.Second):
		t.Fatal("receiveBoth did not return")
	}
	slices.Sort(got)

	want := []string{"data 1", "data 2"}
	if !slices.Equal(got, want) {
		t.Fatalf("receiveBoth returned %v; want %v in any order", got, want)
	}
}
