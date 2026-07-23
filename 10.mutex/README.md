# Lesson 10: Mutex

Previous: [Fan-out/fan-in](../9.fan-in-out/README.md) ·
[Course index](../README.md) ·
Next: [Confinement](../11.confinement/README.md)

## Concept

`sync.Mutex` protects an invariant over shared mutable state. Only one
goroutine can hold the mutex at a time; other callers of `Lock` wait.

In this example, multiplication uses goroutine-local values and needs no lock.
Appending changes the shared slice header and backing storage, so the append is
the critical section.

## Critical-section design

```go
lock.Lock()
defer lock.Unlock()
result = append(result, value)
```

Keep the critical section as small as correctness allows. Expensive independent
work should happen before acquiring the lock so it can proceed concurrently.

`defer` makes the unlock follow every normal return path. For very hot,
carefully measured sections, explicit unlock may reduce overhead, but clarity
and correctness come first.

A deferred call runs when the surrounding function returns, not when the
nearest lexical block ends. It is appropriate in `processNum` because the
function returns immediately after `append`. In a long loop or a function that
does more work afterward, use a small helper function or an explicit unlock so
the critical section does not accidentally grow.

## Safety is not ordering

The mutex prevents concurrent mutation of the slice. It does not decide which
goroutine acquires the lock first, so the result can be any permutation of the
doubled numbers.

`wg.Wait` ensures all appends finish before `main` reads the slice.

## Common mistakes

- Protecting writes but leaving concurrent reads unlocked
- Copying a mutex after it has been used
- Calling code with unknown blocking behavior while holding a lock
- Returning without unlocking
- Assuming a mutex makes operations fair or ordered
- Using one global lock for unrelated state

## Mutex or channel?

Use a mutex when several operations must synchronously access the same state.
Use channels when transferring ownership, coordinating lifecycles, or modeling
a stream of work. They are complementary, not competing rules.

## Run it

```bash
go run ./10.mutex
go run -race ./10.mutex
```

## Experiments

- Remove the mutex and use the race detector.
- Preallocate the result and give each worker a unique index, as lesson 11 does.
- Move multiplication under the lock and explain the lost concurrency.

## References

- [`sync.Mutex` documentation](https://pkg.go.dev/sync#Mutex)
- [Go memory model: locks](https://go.dev/ref/mem)
- [Go race detector](https://go.dev/doc/articles/race_detector)
