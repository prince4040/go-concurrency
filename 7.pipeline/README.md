# Lesson 7: Pipeline

Previous: [Done channel](../6.done-channel/README.md) ·
[Course index](../README.md) · Next: [Generator](../8.generator/README.md)

## Concept

A pipeline is a series of stages connected by channels. A stage normally:

1. Receives values from an inbound channel.
2. Transforms or filters them.
3. Sends values on an outbound channel.
4. Closes its outbound channel when it stops.

This example converts a slice to a stream, squares each value, and consumes the
result.

## Data flow

```text
[]int → sliceToChannel → <-chan int → sq → <-chan int → main
```

Each stage starts one goroutine and immediately returns its output channel.
That allows stages to be composed before all the data has been processed.
Unbuffered channels provide backpressure: a stage cannot outrun the next stage.

## Ownership and cancellation

`sliceToChannel` owns the first output; `sq` owns the second. A stage never
closes its input because another component owns it.

The stages observe `done` around each operation that can block:

- The source selects between cancellation and sending.
- The square stage selects between cancellation and receiving, then between
  cancellation and sending.

This matters when a consumer stops early. Cancellation must be able to unblock
every upstream stage.

## Common mistake: consuming only part of a pipeline

```go
result := sq(done, input)
fmt.Println(<-result)
return
```

Without closing `done`, `sq` can block sending its next result and the source
can block behind it. Returning from a short command hides the leaked
goroutines; a long-running process accumulates them.

## Guarantees

- Each input value produces one squared value while the pipeline is not
  canceled.
- Values retain source order because each stage has one worker.
- Closure propagates downstream.
- Cancellation may discard work already computed but not delivered.

## Run it

```bash
go run ./7.pipeline
```

## Experiments

- Add another `sq` stage and predict the output.
- Stop after two results, close `done`, and wait for the final channel to close.
- Buffer one stage's output and discuss the memory/backpressure tradeoff.

## References

- [Go concurrency patterns: pipelines and cancellation](https://go.dev/blog/pipelines)
- [Go memory model](https://go.dev/ref/mem)
