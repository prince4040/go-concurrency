package main

/*
A channel communicates typed values between goroutines and also synchronizes
the matching send and receive. On an unbuffered channel, the sender waits until
a receiver is ready.

The goroutine that produces values owns the responsibility for closing its
outbound channel. A receive returns ok=false only after the channel is closed
and all previously sent values have been received.
*/

import "fmt"

func send(out chan<- string) {
	defer close(out)
	out <- "data"
}

func main() {
	c := make(chan string)

	go send(c)

	result, open := <-c
	fmt.Println(result, open)

	_, open = <-c
	fmt.Println("channel open:", open)
	fmt.Println("FINISH!")
}
