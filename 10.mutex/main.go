package main

import (
	"fmt"
	"sync"
)

var lock sync.Mutex

func processNum(wg *sync.WaitGroup, num int, dest *[]int) {
	defer wg.Done()
	processedNum := num * 2

	lock.Lock()
	*dest = append(*dest, processedNum)
	lock.Unlock()
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	result := []int{}

	var wg sync.WaitGroup

	for _, num := range nums {
		wg.Add(1)
		go processNum(&wg, num, &result)
	}

	wg.Wait()
	fmt.Println(result)
}
