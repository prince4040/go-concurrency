# Lesson 1: Goroutines

[Course index](../README.md) · Next: [WaitGroup](../2.waitGroup/README.md)

## Concept

A goroutine is a function execution managed by the Go runtime. The `go`
keyword starts the call independently and lets the caller continue immediately.
Goroutines share the process address space, so they are inexpensive to create
but still require synchronization when they access the same data.

## Mental model

`main` is also a goroutine. The program ends when `main` returns; Go does not
implicitly wait for other goroutines.

In this lesson:

1. `main` starts three calls to `printNumber`.
2. The runtime chooses when each call runs.
3. `wg.Add(1)` registers each task before its `go` statement starts it.
4. Each goroutine defers `wg.Done()` so completion is always recorded when the
   function returns.
5. A WaitGroup keeps `main` from returning early.
6. `FINISH!` prints only after all three calls return.

The three numbers may appear in any order.

## Ownership and lifecycle

The goroutines only read their own `num` value. This example calls
`fmt.Println`, which writes through concurrency-safe `os.Stdout`; do not
generalize that property to every `io.Writer` passed to `fmt.Fprintf`. The
WaitGroup owns no work; it only records task completion.

Go 1.22 and later create new iteration variables for each `for range`
iteration, so each closure receives the intended `num`. Older Go versions
required an explicit per-iteration copy. A variable declared outside the loop
can still be shared by every closure and needs the usual synchronization.

## Common mistake: sleeping instead of synchronizing

```go
go printNumber("1")
time.Sleep(time.Second)
```

This guesses how long the goroutine needs. A loaded machine, slow I/O, or a
future code change can invalidate the guess. Use a WaitGroup, channel, or
context that expresses the actual completion condition.

## Run it

```bash
go run ./1.routine
```

Run it several times. The ordering can change, but all numbers must appear
before `FINISH!`.

## Experiments

- Increase the number of goroutines to 100.
- Print a message immediately after each `go` statement and compare scheduling.
- Remove `wg.Wait` temporarily and observe why some work may disappear.

## When to use goroutines

Use goroutines when operations can make independent progress, such as handling
requests, waiting for I/O, or processing independent work. Do not add a
goroutine merely to make a function “faster”; coordination has a cost and the
work may not be parallelizable.

## References

- [Go specification: go statements](https://go.dev/ref/spec#Go_statements)
- [Go memory model: goroutine creation](https://go.dev/ref/mem)
