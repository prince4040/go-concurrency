package main

import "fmt"

/*
A for/select loop multiplexes channels repeatedly rather than choosing only one
operation. Receiving from a closed channel is always ready, so a closed input
must be disabled to keep it from dominating the select.

Assigning nil to a local channel variable disables that select case because
communication on a nil channel can never proceed. The loop finishes only after
both producers have closed their own output channels.
*/

func sendValues(out chan<- string, values ...string) {
	defer close(out)

	for _, value := range values {
		out <- value
	}
}

func consumeStreams(
	letters, numbers <-chan string,
	consume func(kind, value string),
) {
	for letters != nil || numbers != nil {
		select {
		case value, open := <-letters:
			if !open {
				letters = nil
				continue
			}
			consume("letter", value)
		case value, open := <-numbers:
			if !open {
				numbers = nil
				continue
			}
			consume("number", value)
		}
	}
}

func main() {
	letters := make(chan string)
	numbers := make(chan string)

	go sendValues(letters, "a", "b", "c")
	go sendValues(numbers, "1", "2", "3")

	consumeStreams(letters, numbers, func(kind, value string) {
		fmt.Printf("%s: %s\n", kind, value)
	})
}
