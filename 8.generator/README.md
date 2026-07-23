# Lesson 8: Generator

Previous: [Pipeline](../7.pipeline/README.md) ·
[Course index](../README.md) ·
Next: [Fan-out/fan-in](../9.fan-in-out/README.md)

## Concept

A generator turns a function into a stream. `repeatFn` calls a function
repeatedly and emits its results until cancellation. `take` limits any stream
to at most `n` values.

The composition:

```go
for value := range take(done, repeatFn(done, random), 10) {
	fmt.Println(value)
}
```

lets the consumer request a finite prefix from an otherwise infinite source.

## Generic stream type

`repeatFn[T any]` and `take[T any]` use a type parameter. At each call, `T`
stands for one concrete value type, such as `int`, and that same type flows
through the input and output channels. Generics let the concurrency pattern be
reused without changing its lifecycle behavior.

## Why `take` has two selects

Receiving and forwarding are separate blocking operations:

1. Wait for cancellation or the next input value.
2. Stop if the input closed.
3. Wait for cancellation or a downstream receiver.

A tempting shortcut is incorrect:

```go
select {
case <-done:
	return
case out <- <-input:
}
```

The Go specification requires send values in select cases to be evaluated when
entering the select. The inner `<-input` can therefore block before `done` gets
a chance to win.

## Closure behavior

`take` closes its output when:

- it has forwarded `n` values,
- `done` is closed, or
- its input closes early.

Checking the comma-ok result prevents a closed input from generating an
unlimited sequence of zero values.

## Generator limitation

The preliminary cancellation check avoids starting another `fn` call after
cancellation has already been observed. It cannot give cancellation priority:
`done` may close immediately after the check. Cancellation also cannot
interrupt `fn` while `fn` itself is running. A generator function should return
promptly or accept its own cancellation mechanism if it performs blocking work.

## Run it

```bash
go run ./8.generator
```

The random values differ on every run; the guaranteed property is that at most
ten are printed.

## Experiments

- Replace the random function with a counter closure.
- Close the input before `n` values and observe `take` close early.
- Stop reading from `take`, then close `done` and confirm its output closes.

## Tests

```bash
go test ./8.generator
go test -race ./8.generator
```

The tests cover the exact limit, early input closure, cancellation while
receiving, and cancellation while forwarding.

## References

- [Go specification: select evaluation](https://go.dev/ref/spec#Select_statements)
- [Go concurrency patterns: pipelines](https://go.dev/blog/pipelines)
