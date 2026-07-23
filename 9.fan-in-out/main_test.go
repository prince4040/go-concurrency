package main

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestIsPrime(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want bool
	}{
		{name: "negative", num: -3, want: false},
		{name: "zero", num: 0, want: false},
		{name: "one", num: 1, want: false},
		{name: "two", num: 2, want: true},
		{name: "composite", num: 25, want: false},
		{name: "prime", num: 29, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isPrime(test.num); got != test.want {
				t.Fatalf("isPrime(%d) = %v; want %v", test.num, got, test.want)
			}
		})
	}
}

func TestFanInForwardsEveryValueAndCloses(t *testing.T) {
	done := make(chan struct{})
	first := bufferedStream(1, 3, 5)
	second := bufferedStream(2, 4, 6)

	got := collectInts(t, fanIn(done, first, second))
	slices.Sort(got)

	want := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fanIn returned %v; want %v", got, want)
	}
}

func TestPrimeFinderFiltersValuesAndCloses(t *testing.T) {
	done := make(chan struct{})
	input := bufferedStream(-1, 0, 1, 2, 3, 4, 5, 9)

	got := collectInts(t, primeFinder(done, input))
	want := []int{2, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("primeFinder returned %v; want %v", got, want)
	}
}

func TestPrimeFinderStopsWhenCanceledWhileInputIsIdle(t *testing.T) {
	done := make(chan struct{})
	idle := make(chan int)
	output := primeFinder(done, idle)

	close(done)
	waitUntilClosed(t, output)
}

func TestFanInStopsWhenCanceledWhileInputIsIdle(t *testing.T) {
	done := make(chan struct{})
	idle := make(chan int)
	output := fanIn(done, idle)

	close(done)
	waitUntilClosed(t, output)
}

func TestFanInStopsWhenCanceledWhileOutputIsBlocked(t *testing.T) {
	done := make(chan struct{})
	input := make(chan int)
	_, completed := fanInWithCompletion(done, input)
	sent := make(chan struct{})

	go func() {
		input <- 42
		close(sent)
	}()

	waitForTestSignal(t, sent, "fanIn did not receive its input value")
	close(done)
	waitForTestSignal(t, completed, "fanIn did not stop after cancellation")
}

func bufferedStream(values ...int) <-chan int {
	stream := make(chan int, len(values))
	for _, value := range values {
		stream <- value
	}
	close(stream)
	return stream
}

func waitUntilClosed[T any](t *testing.T, channel <-chan T) {
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

func collectInts(t *testing.T, channel <-chan int) []int {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	var values []int
	for {
		select {
		case value, open := <-channel:
			if !open {
				return values
			}
			values = append(values, value)
		case <-timer.C:
			t.Fatal("channel did not close")
			return nil
		}
	}
}

func waitForTestSignal(t *testing.T, signal <-chan struct{}, failure string) {
	t.Helper()

	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal(failure)
	}
}
