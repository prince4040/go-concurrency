package main

import "fmt"

/*
A closed done channel is a broadcast cancellation signal: every goroutine
waiting to receive from it becomes ready. The empty struct communicates no
payload and makes the intent explicit.

The producer owns and closes its output channel. The caller owns and closes
done, then waits for the output to close so main cannot exit while cleanup is
still in progress.
*/

func fibonacci(done <-chan struct{}) <-chan int {
	numbers := make(chan int)

	go func() {
		defer close(numbers)

		x, y := 0, 1
		for {
			select {
			case <-done:
				return
			case numbers <- x:
				x, y = y, x+y
			}
		}
	}()

	return numbers
}

func main() {
	done := make(chan struct{})
	numbers := fibonacci(done)

	for range 10 {
		fmt.Println(<-numbers)
	}

	close(done)

	// Ranging until close is also an explicit wait for producer shutdown.
	for range numbers {
	}

	fmt.Println("FINISH!")
}
