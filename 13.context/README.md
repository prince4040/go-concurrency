# Lesson 13: Context

Previous: [Or-done channel](../12.or-done-channel/README.md) ·
[Course index](../README.md)

## Concept

`context.Context` carries cancellation, deadlines, and request-scoped values
through a call tree. It standardizes the done-channel pattern for operations
that cross function and API boundaries.

The example creates:

- a parent context with a 400 ms deadline,
- a child context with a 200 ms deadline,
- an apple producer governed by the child, and
- an orange producer governed by the parent.

The child stops first. Its cancellation does not stop the parent. When the
parent deadline expires, the remaining producer closes its output and both
consumers finish.

## Cancellation hierarchy

```text
Background
└── parent: 400 ms
    ├── child: 200 ms → apple producer
    └── orange producer
```

Canceling `parent` would cancel both branches. Canceling `child` affects only
the child branch.

## API rules

- Pass context explicitly as the first parameter, conventionally named `ctx`.
- Never pass a nil context; use `context.TODO()` when the correct parent is not
  yet known.
- Call every returned cancel function, even when a deadline is expected to
  expire, so timers and parent references are released promptly.
- Do not store context in a struct unless an API specifically requires it.
- Use context values only for request-scoped data crossing API boundaries, not
  as optional function arguments.

## Context does not wait

Calling a cancel function signals cancellation and arranges for `ctx.Done` to
close. The `Context` contract permits that closure to happen asynchronously
after the cancel function returns. Cancellation never waits for goroutines to
return, so callers must observe `Done` or worker completion explicitly. This
example combines context with:

- producers that close their own output channels,
- consumers that range until closure, and
- a WaitGroup that observes consumer completion.

## Common mistake: circular shutdown

```go
defer cancel()
wg.Wait() // worker waits for ctx.Done
```

The deferred cancel cannot run until after `Wait`, while the worker cannot
finish until cancellation. Use an explicit cancellation event, a deadline, or
a coordinator that cancels before waiting.

## Timing and correctness

The ticker simulates periodic production; it does not provide synchronization
for shutdown. Context cancellation and channel closure provide the lifecycle
guarantees. Exact message counts near a deadline are intentionally
nondeterministic.

## Run it

```bash
go run ./13.context
```

## Experiments

- Cancel the parent explicitly before its deadline.
- Give the child a deadline later than the parent and observe inherited
  cancellation.
- Print `context.Cause(ctx)` and compare it with `ctx.Err()`.

## Tests

```bash
go test ./13.context
go test -race ./13.context
```

The tests verify prompt producer closure after cancellation and bounded
completion of the complete example.

## References

- [`context` package documentation](https://pkg.go.dev/context)
- [Go concurrency patterns: context](https://go.dev/blog/context)
- [Canceling in-progress operations](https://go.dev/doc/database/cancel-operations)
