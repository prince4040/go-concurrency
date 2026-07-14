package main

import (
	"fmt"
	"time"
)

func someFun(c chan string) {
	arr := []string{"a", "b", "c", "d"}

	for _, char := range arr {
		select {
		case c <- char:
		}
	}

	close(c)
}

func main() {
	c1 := make(chan string, 3)

	go someFun(c1)

	for v := range c1 {
		fmt.Println(v)
	}

	time.Sleep(10 * time.Second)
}
