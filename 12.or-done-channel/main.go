package main

import (
	"fmt"
)

/*
orDone adapts a channel so a consumer can range over it while still responding
to cancellation. Its forwarding goroutine observes done while waiting for
input and again while waiting for the downstream consumer.

Both checks matter. Omitting either can strand the forwarding goroutine when an
upstream producer is idle or a downstream consumer stops receiving.
*/

func producer(done <-chan struct{}) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case out <- "data":
			}
		}
	}()

	return out
}

func orDone[T any](done <-chan struct{}, c <-chan T) <-chan T {
	out, _ := orDoneWithCompletion(done, c)
	return out
}

// orDoneWithCompletion makes forwarding-goroutine termination observable.
func orDoneWithCompletion[T any](
	done <-chan struct{},
	c <-chan T,
) (<-chan T, <-chan struct{}) {
	out := make(chan T)
	completed := make(chan struct{})

	go func() {
		defer close(completed)
		defer close(out)

		for {
			select {
			case <-done:
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- val:
				}
			}
		}
	}()

	return out, completed
}

func main() {
	done := make(chan struct{})
	stream := producer(done)
	values, forwarderCompleted := orDoneWithCompletion(done, stream)

	for range 5 {
		fmt.Println(<-values)
	}

	close(done)

	<-forwarderCompleted
	for range stream {
	}
}
