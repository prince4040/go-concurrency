package main

import (
	"fmt"
	"sync"
)

/*
A goroutine is an independently scheduled function execution in the same
process. Starting a goroutine does not make the caller wait for it, and the
program ends as soon as main returns.

The order of the printed numbers is intentionally unspecified: scheduling is
a runtime decision. The explicit Add, go, and Done sequence keeps the
goroutine primitive visible; lesson 2 explains the WaitGroup coordination in
detail and introduces its modern Go helper.
*/

func printNumber(num string) {
	fmt.Println(num)
}

func runNumbers(nums []string, print func(string)) {
	var wg sync.WaitGroup

	for _, num := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			print(num)
		}()
	}

	wg.Wait()
}

func main() {
	runNumbers([]string{"1", "2", "3"}, printNumber)
	fmt.Println("FINISH!")
}
