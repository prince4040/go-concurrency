package main

import (
	"fmt"
	"math/rand/v2"
	"time"
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

func primeFinder(done <-chan bool, intStream <-chan int) <-chan int {
	primes := make(chan int)

	isPrime := func(num int) bool {
		for i := 2; i < num; i++ {
			if num%i == 0 {
				return false
			}
		}
		return true
	}

	go func() {
		defer close(primes)

		for {
			select {
			case <-done:
				return
			case num := <-intStream:
				if isPrime(num) {
					primes <- num
				}
			}
		}
	}()

	return primes
}

func getRandomNum() int {
	return rand.IntN(5000000000)
}

func main() {
	start := time.Now()

	done := make(chan bool)
	defer close(done)

	randNumStream := repeatFun(done, getRandomNum)
	primeStream := primeFinder(done, randNumStream)

	for v := range take(done, primeStream, 10) {
		fmt.Println(v)
	}

	fmt.Println(time.Since(start))
}
