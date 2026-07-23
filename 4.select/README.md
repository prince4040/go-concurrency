# Lesson 4: Select

Previous: [Channels](../3.channels/README.md) ·
[Course index](../README.md) · Next: [For/select](../5.for-select/README.md)

## Concept

`select` waits on multiple channel sends or receives. When one operation can
proceed, its case runs. With no `default`, the whole statement blocks until at
least one case is ready.

If several cases are ready together, Go makes a uniform pseudo-random choice.
The first case in source order has no priority.

## Walkthrough

Two producers send on separate channels with capacity one and close their own
outputs:

1. `select` receives whichever result becomes ready first.
2. That selected case runs once.
3. A normal receive drains the other channel.
4. Both buffered producers can finish even before a consumer chooses them.

The two lines may print in either order.

## Common mistake: abandoning the losing sender

```go
select {
case value := <-first:
	fmt.Println(value)
case value := <-second:
	fmt.Println(value)
}
// main returns without receiving the other send
```

With unbuffered outputs, the unselected producer remains blocked. A short
command then exits, hiding the problem; the same pattern in a server leaks a
goroutine. Receive the remaining work, buffer deliberately as this example
does, or cancel the unchosen operation.

## `default` is not a faster select

A `default` case runs immediately when no communication is ready. Repeating
that select in a tight loop can consume a CPU core. Use `default` only when
non-blocking behavior is genuinely required, and arrange an appropriate wait
or backoff.

## Run it

```bash
go run ./4.select
```

## Experiments

- Run the example repeatedly and record which channel wins first.
- Change both channels to unbuffered channels and explain why draining the
  second result remains necessary.
- Add a timeout case with `time.After`.
- Add a third input and select one result from the three.

Lesson 5 introduces the state needed to select repeatedly until several inputs
close.

## References

- [Go specification: select statements](https://go.dev/ref/spec#Select_statements)
- [Go specification: receive operator](https://go.dev/ref/spec#Receive_operator)
