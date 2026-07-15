package main

import (
	"fmt"
	"sync"
	"time"
)

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
}
