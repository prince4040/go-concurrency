# Lesson 2: WaitGroup

Previous: [Goroutines](../1.routine/README.md) ·
[Course index](../README.md) · Next: [Channels](../3.channels/README.md)

## Concept

`sync.WaitGroup` waits for a collection of tasks. In Go 1.25 and later,
`WaitGroup.Go` is the preferred way to start and track a goroutine together:

```go
wg.Go(func() {
	runJob(id)
})
```

`Wait` blocks until every tracked function returns. It controls completion, not
execution order.

Lesson 1 wrote the underlying `Add`, `go`, and `Done` sequence explicitly so
the goroutine-starting statement remained visible. `WaitGroup.Go` packages that
common sequence into one operation and prevents mismatched task counts.

## Mental model

Think of the WaitGroup as a counter:

- Starting a tracked task increments it.
- Returning from that task decrements it.
- `Wait` unblocks when the count reaches zero.

The lesson starts three jobs. Their messages are unordered, but the final
message is ordered after all three by `wg.Wait`.

## Classic Add/Done form

Older code and code that tracks work not started by `WaitGroup.Go` uses:

```go
wg.Add(1)
go func() {
	defer wg.Done()
	runJob(id)
}()
```

The positive `Add` must happen before the goroutine starts and before a `Wait`
that could see a zero count. Calling `Add` inside the goroutine creates a race
between `Add` and `Wait`.

## Correctness rules

- Do not copy a WaitGroup after first use; pass a pointer when a function needs
  to access it.
- Do not let the counter become negative.
- Do not start a new independent batch until the previous `Wait` has returned.
- A function passed to `WaitGroup.Go` must not panic.
- A WaitGroup does not collect results or errors; use channels or a higher-level
  abstraction when tasks must return them.

## Run it

```bash
go run ./2.waitGroup
```

## Experiments

- Increase the job count and confirm the final line remains last.
- Rewrite the loop using `Add` and `defer Done`.
- Pass results through a channel and compare completion coordination with data
  communication.

## When to use it

Use a WaitGroup when one scope starts a known set of concurrent tasks and needs
all of them to finish. It is not cancellation: `Wait` cannot tell a task to
stop.

## References

- [`sync.WaitGroup` documentation](https://pkg.go.dev/sync#WaitGroup)
- [Go memory model](https://go.dev/ref/mem)
