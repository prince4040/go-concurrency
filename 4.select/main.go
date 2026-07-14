package main

/*
	select lets goroutine wait on multiple comminicatuion opearations
*/

import (
	"fmt"
)

func someFunc(n string, c chan string) {
	c <- "data " + n

	// close(c)
	// if n == "2" {
	// 	time.Sleep(5 * time.Second)
	// }
}

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go someFunc("1", c1)
	go someFunc("2", c2)

	// only runs one block which can be executed first
	select {
	case result := <-c1:
		fmt.Println(result)
	case result := <-c2:
		fmt.Println(result)
	}

	fmt.Println("FINISH!")
}
