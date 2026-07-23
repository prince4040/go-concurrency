package main

import "fmt"

/*
A pipeline is a series of stages connected by channels. Each stage receives
values, transforms them, and owns the channel it returns. Closing an outbound
channel tells the next stage that no more values will arrive.

Every blocking receive and send also observes done. That lets cancellation
propagate upstream when a downstream consumer stops before draining the
pipeline.
*/

func sliceToChannel(done <-chan struct{}, nums []int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, num := range nums {
			select {
			case <-done:
				return
			case out <- num:
			}
		}
	}()

	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for {
			var (
				num  int
				open bool
			)

			select {
			case <-done:
				return
			case num, open = <-in:
				if !open {
					return
				}
			}

			select {
			case <-done:
				return
			case out <- num * num:
			}
		}
	}()

	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	nums := []int{1, 2, 3, 4, 5, 6}

	dataChannel := sliceToChannel(done, nums)
	finalChannel := sq(done, dataChannel)

	for num := range finalChannel {
		fmt.Println(num)
	}
}
