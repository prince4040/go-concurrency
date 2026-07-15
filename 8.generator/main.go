package main

import (
	"fmt"
	"math/rand/v2"
)

// Generator
func repeatFun[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()

	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)

		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()

	return taken
}

func main() {
	done := make(chan bool)
	defer close(done)

	getRandomNum := func() int {
		return rand.IntN(50000000)
	}

	randNumStream := repeatFun(done, getRandomNum)

	for v := range take(done, randNumStream, 10) {
		fmt.Println(v)
	}
}
