package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(ctx context.Context, msg string, c chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case c <- msg:
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appleChannel := make(chan string)
	orangeChannel := make(chan string)

	go producer(ctx, "apple", appleChannel)
	go producer(ctx, "orange", orangeChannel)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func1(ctx, &wg, appleChannel)

	wg.Add(1)
	go func2(ctx, &wg, orangeChannel)

	wg.Wait()
}

func orDone[T any](ctx context.Context, c <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()

	return out
}

func func1(ctx context.Context, parentWg *sync.WaitGroup, c <-chan string) {
	defer parentWg.Done()

	doWork := func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for v := range orDone(ctx, c) {
			fmt.Println(v)
		}
	}

	wg := sync.WaitGroup{}
	newCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	for range 3 {
		wg.Add(1)
		go doWork(newCtx, &wg)
	}

	wg.Wait()
}

func func2(ctx context.Context, parentWg *sync.WaitGroup, c <-chan string) {
	defer parentWg.Done()

	for v := range orDone(ctx, c) {
		fmt.Println(v)
	}
}
