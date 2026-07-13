package main

import (
	"fmt"
	"time"
)

func someFun(num string) {
	fmt.Println(num)
}

func main() {
	go someFun("1")
	go someFun("2")
	go someFun("3")

	time.Sleep(1 * time.Second)

	fmt.Println("FINISH!")
}
