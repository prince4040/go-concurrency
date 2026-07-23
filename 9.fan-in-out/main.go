package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"
)

/*
Fan-out starts multiple workers that compete to receive from one input,
distributing work rather than copying every value to every worker. Fan-in
merges their output channels into one stream.

The number of printed results is deterministic here, but their values, order,
and which worker handles each value are not. More workers create parallelism
only when the runtime and workload can use it; concurrency alone does not
guarantee a speedup.
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
	taken := make(chan T)

	go func() {
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

	return taken
}

func isPrime(num int) bool {
	if num < 2 {
		return false
	}

	for divisor := 2; divisor <= num/divisor; divisor++ {
		if num%divisor == 0 {
			return false
		}
	}

	return true
}

func primeFinder(done <-chan struct{}, intStream <-chan int) <-chan int {
	primes := make(chan int)

	go func() {
		defer close(primes)

		for {
			var (
				num  int
				open bool
			)

			select {
			case <-done:
				return
			case num, open = <-intStream:
				if !open {
					return
				}
			}

			if !isPrime(num) {
				continue
			}

			select {
			case <-done:
				return
			case primes <- num:
			}
		}
	}()

	return primes
}

func getRandomNum() int {
	return rand.IntN(50_000_000)
}

func fanIn[T any](done <-chan struct{}, channels ...<-chan T) <-chan T {
	output, _ := fanInWithCompletion(done, channels...)
	return output
}

// fanInWithCompletion makes the merge goroutines' termination observable.
func fanInWithCompletion[T any](
	done <-chan struct{},
	channels ...<-chan T,
) (<-chan T, <-chan struct{}) {
	fannedInChannel := make(chan T)
	completed := make(chan struct{})
	var wg sync.WaitGroup

	transfer := func(c <-chan T) {
		for {
			var (
				value T
				open  bool
			)

			select {
			case <-done:
				return
			case value, open = <-c:
				if !open {
					return
				}
			}

			select {
			case <-done:
				return
			case fannedInChannel <- value:
			}
		}
	}

	for _, channel := range channels {
		wg.Go(func() {
			transfer(channel)
		})
	}

	go func() {
		defer close(completed)
		wg.Wait()
		close(fannedInChannel)
	}()

	return fannedInChannel, completed
}

func main() {
	start := time.Now()

	done := make(chan struct{})

	randNumStream := repeatFn(done, getRandomNum)

	numWorkers := runtime.GOMAXPROCS(0)
	primeFinderChannels := make([]<-chan int, numWorkers)

	for i := range numWorkers {
		primeFinderChannels[i] = primeFinder(done, randNumStream)
	}

	fannedInPrimeStream, fanInCompleted := fanInWithCompletion(
		done,
		primeFinderChannels...,
	)

	for num := range take(done, fannedInPrimeStream, 10) {
		fmt.Println(num)
	}

	close(done)
	<-fanInCompleted
	for _, channel := range primeFinderChannels {
		for range channel {
		}
	}
	for range randNumStream {
	}

	fmt.Println(time.Since(start))
}
