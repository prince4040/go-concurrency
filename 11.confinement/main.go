package main

import (
	"fmt"
	"sync"

	"github.com/prince4040/concurrency/11.confinement/balance"
)

/*
This lesson shows two related ownership techniques.

The first partitions a preallocated result slice. Every goroutine receives a
unique element, so no two goroutines write the same memory location. Wait waits
for all writes before main reads the completed slice.

The balance example uses stronger channel confinement: one account goroutine
is the only code allowed to access the balance. Other goroutines communicate
requests instead of sharing that state directly.
*/

func processNum(num int, dest *int) {
	*dest = num * 2
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	result := make([]int, len(nums))

	var wg sync.WaitGroup

	for i, num := range nums {
		wg.Go(func() {
			processNum(num, &result[i])
		})
	}

	wg.Wait()
	fmt.Println(result)

	balance.BalanceMain()
}
