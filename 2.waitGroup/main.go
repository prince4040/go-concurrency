package main

import (
	"fmt"
	"sync"
)

/*
A WaitGroup represents a set of tasks that must finish before execution can
continue. WaitGroup.Go starts a goroutine and tracks its lifetime as one
operation, avoiding mismatched Add and Done calls.

Wait only establishes completion; it does not prescribe task order. The jobs
below can print in any order, but "all jobs complete" is guaranteed to print
after all of them.
*/

func runJob(id int) {
	fmt.Printf("job %d complete\n", id)
}

func runJobs(ids []int, run func(int)) {
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Go(func() {
			run(id)
		})
	}

	wg.Wait()
}

func main() {
	runJobs([]int{1, 2, 3}, runJob)
	fmt.Println("all jobs complete")
}
