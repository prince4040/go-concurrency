# Lesson 5: For/select

Previous: [Select](../4.select/README.md) ·
[Course index](../README.md) ·
Next: [Done channel](../6.done-channel/README.md)

## Concept

A `for` around `select` repeatedly multiplexes independent streams. This is
common in event loops, coordinators, and components that combine results.
Unlike Lesson 4's single selection followed by an explicit drain, this lesson
does not know which input will close first or how many values each will send.

The difficult part is termination: a receive from a closed channel is always
ready and immediately returns the zero value.

## Mental model

The two producers own and close their output channels. The consumer tracks
each channel in a local variable:

1. Select a value from either input.
2. When a receive reports `open=false`, assign `nil` to that variable.
3. A nil channel can never communicate, so its select case is disabled.
4. Exit when both variables are nil.

This prevents a closed channel from repeatedly winning the select.

## Common mistake: ignoring comma-ok

```go
for {
	select {
	case value := <-input:
		fmt.Println(value)
	}
}
```

After `input` closes, this loop prints zero values forever. Always decide how a
closed input changes the event loop's state.

## Why the ordering changes

Each producer preserves its own order: `a` comes before `b`, and `1` comes
before `2`. The interleaving between the streams is not defined because select
chooses among the currently ready operations.

## Run it

```bash
go run ./5.for-select
```

## Experiments

- Give one channel a buffer and compare interleavings.
- Add a third producer and extend the loop condition.
- Temporarily omit one `channel = nil` assignment and observe the failure.

## When to use it

Use for/select when one goroutine coordinates several ongoing channel events.
If only one channel is involved, a simple `for value := range channel` is
usually clearer.

## References

- [Go specification: select](https://go.dev/ref/spec#Select_statements)
- [Go specification: receive operator](https://go.dev/ref/spec#Receive_operator)
