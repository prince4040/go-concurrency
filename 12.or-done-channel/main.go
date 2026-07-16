package main

import (
	"fmt"
	"sync"
)

func producer(done <-chan any, c chan<- string) {
	for {
		select {
		case <-done:
			return
		case c <- "data":
		}
	}
}

func receiver(wg *sync.WaitGroup, done <-chan any, c <-chan string) {
	defer wg.Done()

	for val := range orDone(done, c) {
		// complex operation
		fmt.Println(val)
	}

	/*
		NAIVE
		for {
			select {
			case <-done:
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				// complex operation
				fmt.Println(val)

			}
		}
	*/
}

func orDone[T any](done <-chan any, c <-chan T) <-chan T {
	out := make(chan T)

	go func() {
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

	return out
}

func main() {
	done := make(chan any)
	defer close(done)
	c := make(chan string)

	wg := sync.WaitGroup{}

	go producer(done, c)

	wg.Add(1)
	go receiver(&wg, done, c)

	wg.Wait()
}
