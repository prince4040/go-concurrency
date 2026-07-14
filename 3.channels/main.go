package main

/*
	channels are the way for communication between multiple routines
*/

import "fmt"

func someFunc(_ int, c chan string) {
	c <- "data"

	// close(c)
}

func main() {
	c := make(chan string)

	go someFunc(1, c)

	/* Blocking: blocks the main routine until channel closes or gets any data from channel */
	result, ok := <-c

	fmt.Println(result, ok)
	fmt.Println("FINISH!")
}
