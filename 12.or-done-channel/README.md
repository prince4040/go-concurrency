# Lesson 12: Or-done Channel

Previous: [Confinement](../11.confinement/README.md) ·
[Course index](../README.md) · Next: [Context](../13.context/README.md)

## Concept

Ranging over a channel is concise:

```go
for value := range input {
	use(value)
}
```

But the receive cannot also observe a cancellation channel. `orDone` wraps the
input with a forwarding goroutine and returns an output that closes when either
the input closes or cancellation occurs.

## Why cancellation is checked twice

The forwarding loop contains two independent blocking points:

1. Waiting for an input value
2. Waiting for the downstream consumer to receive that value

Each point selects against `done`. Checking only the receive side still leaks
when downstream stops; checking only the send side still leaks when upstream is
idle.

## Lifecycle

```text
producer ── stream ──> orDone forwarder ── values ──> main
      ↑                     ↑
      └────── done ─────────┘
```

- The producer closes `stream`.
- The `orDone` forwarder closes `values`.
- `main` closes `done`.
- After cancellation, `main` waits for both outputs to close.

One final value may race with cancellation if both forwarding and cancellation
are ready together. Cancellation is a request to stop, not a priority
mechanism.

## Tradeoff

Each `orDone` call creates a goroutine and another channel. Use it when it makes
composition and lifecycle handling clearer; a direct select is cheaper and
often clearer in a small loop.

## Run it

```bash
go run ./12.or-done-channel
```

The program prints five values, closes `done`, and waits for both goroutines to
finish.

## Experiments

- Replace the producer with an input that never sends, then cancel.
- Stop receiving from `values`, then cancel and wait for it to close.
- Close the input normally instead of canceling.

## Tests

```bash
go test ./12.or-done-channel
go test -race ./12.or-done-channel
```

## References

- [Go concurrency patterns: pipelines and cancellation](https://go.dev/blog/pipelines)
- [Go specification: select](https://go.dev/ref/spec#Select_statements)
