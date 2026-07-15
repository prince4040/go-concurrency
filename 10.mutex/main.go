package main

import (
	"fmt"
	"sync"
	"time"
)

var lock sync.Mutex

func getprocessedNum(num int) int {
	time.Sleep(1 * time.Second)
	return num * 2
}

func processNum(wg *sync.WaitGroup, nums *[]int, num int) {
	defer wg.Done()

	processedNum := getprocessedNum(num)

	lock.Lock()
	*nums = append(*nums, processedNum)
	lock.Unlock()
}

func main() {
	start := time.Now()

	nums := []int{1, 2, 3, 4, 5}
	result := []int{}

	var wg sync.WaitGroup

	for _, num := range nums {
		wg.Add(1)
		go processNum(&wg, &result, num)
	}

	wg.Wait()
	fmt.Println(result)

	fmt.Println(time.Since(start))
}
