# Lesson 9: Fan-out and Fan-in

Previous: [Generator](../8.generator/README.md) ·
[Course index](../README.md) · Next: [Mutex](../10.mutex/README.md)

## Concept

Fan-out starts multiple workers that receive from one input channel. Each value
goes to one worker, so the workers distribute the work.

Fan-in starts one transfer goroutine per input and merges all their values into
one output.

```text
                           ┌→ primeFinder ─┐
random number stream ─────┼→ primeFinder ─┼→ fanIn → take(10)
                           └→ primeFinder ─┘
```

Fan-out is not broadcast. To deliver every value to every consumer, a separate
duplication stage is required.

## Worker count

The example uses `runtime.GOMAXPROCS(0)`, the current maximum number of CPUs that
can execute Go code simultaneously, as a reasonable demonstration worker
count. It is not a universal tuning rule.

More workers can be slower when:

- work is too small to offset coordination,
- contention dominates,
- the bottleneck is elsewhere, or
- the runtime cannot execute more workers in parallel.

The printed duration is an observation, not a benchmark.

## Cancellation-safe fan-in

Each transfer goroutine selects around both:

1. receiving from its input, and
2. sending to the merged output.

After all transfer goroutines stop, a separate closer goroutine closes the
merged channel. Closing earlier would race with active senders.

## Common mistake: range cannot cancel an idle input

```go
for value := range input {
	select {
	case <-done:
		return
	case output <- value:
	}
}
```

If `input` stays open but sends nothing, the range receive blocks before the
select. Put cancellation and receive in the same select.

## Ordering and correctness

- Each input value is handled by at most one prime worker.
- `fanIn` forwards every received value once unless cancellation wins.
- Result ordering is intentionally nondeterministic.
- `isPrime` rejects all values below two and tests divisors only up to the
  square root.

Cancellation is observed around channel operations, but it cannot interrupt
`isPrime` while that function is already computing. Long-running CPU work must
check cancellation inside the computation or be divided into smaller units if
prompt shutdown is required.

## Run it

```bash
go run ./9.fan-in-out
```

## Experiments

- Compare one worker with `runtime.GOMAXPROCS(0)` workers.
- Replace primality testing with cheap work and observe the coordination cost.
- Fan in two small deterministic streams and inspect interleaving.

## Tests

```bash
go test ./9.fan-in-out
go test -race ./9.fan-in-out
```

The tests validate prime boundaries, exact fan-in membership, output closure,
and cancellation while either side is blocked.

## References

- [Go concurrency patterns: pipelines and cancellation](https://go.dev/blog/pipelines)
- [`runtime.GOMAXPROCS` documentation](https://pkg.go.dev/runtime#GOMAXPROCS)
- [`sync.WaitGroup` documentation](https://pkg.go.dev/sync#WaitGroup)
