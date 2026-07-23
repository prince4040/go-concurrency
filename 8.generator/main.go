package main

import (
	"fmt"
	"math/rand/v2"
)

/*
A generator converts a function into a stream of values. repeatFn produces an
unbounded stream, while take turns it into a finite stream for its consumer.

take uses separate cancellation-aware receive and send steps. Writing
"out <- <-in" inside a select is unsafe because the inner receive is evaluated
before Go chooses a select case and can therefore hide cancellation.
*/

func repeatFn[T any](done <-chan struct{}, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			default:
			}

			value := fn()

			select {
			case <-done:
				return
			case stream <- value:
			}
		}
	}()

	return stream
}

func take[T any](done <-chan struct{}, stream <-chan T, n int) <-chan T {
	taken, _ := takeWithCompletion(done, stream, n)
	return taken
}

// takeWithCompletion also exposes worker termination to lifecycle coordinators.
func takeWithCompletion[T any](
	done <-chan struct{},
	stream <-chan T,
	n int,
) (<-chan T, <-chan struct{}) {
	taken := make(chan T)
	completed := make(chan struct{})

	go func() {
		defer close(completed)
		defer close(taken)

		for range n {
			var (
				value T
				open  bool
			)

			select {
			case <-done:
				return
			case value, open = <-stream:
				if !open {
					return
				}
			}

			select {
			case <-done:
				return
			case taken <- value:
			}
		}
	}()

	return taken, completed
}

func main() {
	done := make(chan struct{})

	getRandomNum := func() int {
		return rand.IntN(50000000)
	}

	randNumStream := repeatFn(done, getRandomNum)

	for v := range take(done, randNumStream, 10) {
		fmt.Println(v)
	}

	close(done)
	for range randNumStream {
	}
}
