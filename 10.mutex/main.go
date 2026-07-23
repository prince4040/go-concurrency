package main

import (
	"fmt"
	"sync"
)

/*
A mutex serializes access to shared mutable state. The multiplication can run
concurrently, but append changes the shared slice header and therefore belongs
inside the critical section.

The lock makes the append safe, not ordered: whichever goroutine acquires the
mutex first determines the result order. The mutex is kept beside the state it
protects instead of being a package-level global.
*/

func processNum(lock *sync.Mutex, num int, dest *[]int) {
	processedNum := num * 2

	lock.Lock()
	defer lock.Unlock()

	*dest = append(*dest, processedNum)
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	result := []int{}

	var (
		lock sync.Mutex
		wg   sync.WaitGroup
	)

	for _, num := range nums {
		wg.Go(func() {
			processNum(&lock, num, &result)
		})
	}

	wg.Wait()
	fmt.Println(result)
}
