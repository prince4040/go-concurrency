package main

/*
select waits until one of several channel operations can proceed. If multiple
cases are ready, Go chooses one using a uniform pseudo-random selection; source
order is not a priority.

Each producer has a one-value buffer, so it can finish even when its value is
not selected first. This example performs one select, then explicitly drains
the other result. Lesson 5 adds repeated selection and closure tracking.
*/

import "fmt"

func produce(label string) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		out <- "data " + label
	}()

	return out
}

func receiveBoth(first, second <-chan string) []string {
	select {
	case result := <-first:
		return []string{result, <-second}
	case result := <-second:
		return []string{result, <-first}
	}
}

func main() {
	c1 := produce("1")
	c2 := produce("2")

	for _, result := range receiveBoth(c1, c2) {
		fmt.Println(result)
	}

	fmt.Println("FINISH!")
}
