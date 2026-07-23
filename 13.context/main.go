package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Context carries cancellation and deadlines through a call tree. Canceling a
parent cancels every derived child; canceling a child does not cancel its
parent or siblings.

The apple producer has a shorter child deadline, while the orange producer
lives until the parent deadline. Each producer owns and closes its output, so
consumers can use range and WaitGroup.Wait can observe clean completion.
*/

func producer(ctx context.Context, message string, interval time.Duration) <-chan string {
	if interval <= 0 {
		panic("producer interval must be positive")
	}

	out := make(chan string)

	go func() {
		defer close(out)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}

			select {
			case <-ctx.Done():
				return
			case out <- message:
			}
		}
	}()

	return out
}

func consume(name string, values <-chan string) {
	for value := range values {
		fmt.Printf("%s received %s\n", name, value)
	}
}

func runExample(parentDuration, childDuration, producerInterval time.Duration) {
	parent, cancelParent := context.WithTimeout(context.Background(), parentDuration)
	defer cancelParent()

	child, cancelChild := context.WithTimeout(parent, childDuration)
	defer cancelChild()

	apples := producer(child, "apple", producerInterval)
	oranges := producer(parent, "orange", producerInterval)

	var wg sync.WaitGroup
	wg.Go(func() {
		consume("child", apples)
	})
	wg.Go(func() {
		consume("parent", oranges)
	})

	wg.Wait()
	fmt.Println("parent stopped:", parent.Err())
}

func main() {
	runExample(400*time.Millisecond, 200*time.Millisecond, 50*time.Millisecond)
}
