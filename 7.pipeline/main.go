package main

import "fmt"

func sliceToChannel(nums []int) <-chan int {
	out := make(chan int)

	go func() {
		for num := range nums {
			out <- num
		}
		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for num := range in {
			out <- num * num
		}

		close(out)
	}()

	return out
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6}

	// stage 1
	dataChannel := sliceToChannel(nums)

	// stage 2
	finalChannel := sq(dataChannel)

	// stage 3
	for num := range finalChannel {
		fmt.Println(num)
	}
}
