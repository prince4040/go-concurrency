package main

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
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
			case num, ok := <-intStream:
				if !ok {
					return
				}
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

func fanin[T any, K any](done <-chan K, channels ...<-chan T) <-chan T {
	fannedInChannel := make(chan T)
	var wg sync.WaitGroup

	transfer := func(c <-chan T) {
		defer wg.Done()

		for num := range c {
			select {
			case <-done:
				return
			case fannedInChannel <- num:
			}
		}
	}

	for _, channel := range channels {
		wg.Add(1)
		go transfer(channel)
	}

	go func() {
		wg.Wait()
		close(fannedInChannel)
	}()

	return fannedInChannel
}

func main() {
	start := time.Now()

	done := make(chan bool)
	defer close(done)

	randNumStream := repeatFun(done, getRandomNum)
	// primeStream := primeFinder(done, randNumStream)

	// naive
	// for v := range take(done, primeStream, 10) {
	// 	fmt.Println(v)
	// }

	//**************************************************

	// fan-out
	numCPUs := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, numCPUs)

	for i := 0; i < numCPUs; i++ {
		primeFinderChannels[i] = primeFinder(done, randNumStream)
	}

	// fan-in
	fannedInPrimeStream := fanin(done, primeFinderChannels...)

	for num := range take(done, fannedInPrimeStream, 10) {
		fmt.Println(num)
	}

	fmt.Println(time.Since(start))
}
