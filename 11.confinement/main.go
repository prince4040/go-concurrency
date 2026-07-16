package main

import (
	"fmt"
	"sync"
	"time"
)

/*
Instead of protecting shared memory with synchronization primitives
such as a mutex, this program assigns each goroutine exclusive ownership
of the memory it writes to.

Although the `result` slice is shared, every goroutine receives a pointer
to a unique element (`&result[i]`). Since no two goroutines write to the
same memory location, there is no data race and no mutex is required.

This follows Go's concurrency philosophy:
    "Do not communicate by sharing memory; instead, share memory by communicating."

More generally, designing concurrent programs so that each goroutine owns
its data often reduces contention and avoids the serialization introduced
by mutexes, allowing multiple goroutines to make progress in parallel.
*/

func processNum(wg *sync.WaitGroup, num int, dest *int) {
	defer wg.Done()
	processedNum := num * 2

	*dest = processedNum
}

func main() {
	start := time.Now()

	nums := []int{1, 2, 3, 4, 5}
	result := make([]int, len(nums))

	var wg sync.WaitGroup

	for i, num := range nums {
		wg.Add(1)
		go processNum(&wg, num, &result[i])
	}

	wg.Wait()
	fmt.Println(result)

	fmt.Println(time.Since(start))

	// balance.BalanceMain()
}
