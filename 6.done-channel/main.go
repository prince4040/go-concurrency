package main

import (
	"fmt"
)

/*
	to prevent go routines from running infinitely (unintentionally)
*/

func fibonacci(c chan int, done <-chan bool) {
	defer close(c)
	x, y := 0, 1

	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-done:
			fmt.Println("quit")
			return
		}
	}
}

func someFun(done <-chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			fmt.Println("DOING WORK!!")
		}
	}
}

func main() {
	/*
		done := make(chan bool)
		go someFun(done)

		time.Sleep(5 * time.Second)
		close(done)
	*/

	c := make(chan int)
	done := make(chan bool)
	go fibonacci(c, done)

	for range 10 {
		fmt.Println(<-c)
	}
	close(done)

	fmt.Println("FINISH!!!")
}
