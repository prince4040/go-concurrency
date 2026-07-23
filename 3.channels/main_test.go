package main

import (
	"testing"
	"time"
)

func TestSendEmitsValueAndClosesOutput(t *testing.T) {
	output := make(chan string)
	go send(output)

	if got := receiveString(t, output); got != "data" {
		t.Fatalf("send emitted %q; want %q", got, "data")
	}

	select {
	case _, open := <-output:
		if open {
			t.Fatal("output remained open after the final value")
		}
	case <-time.After(time.Second):
		t.Fatal("output did not close")
	}
}

func receiveString(t *testing.T, input <-chan string) string {
	t.Helper()

	select {
	case value, open := <-input:
		if !open {
			t.Fatal("input closed before emitting a value")
		}
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for a value")
		return ""
	}
}
