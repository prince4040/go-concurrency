# Lesson 11: Confinement and Ownership

Previous: [Mutex](../10.mutex/README.md) ·
[Course index](../README.md) ·
Next: [Or-done channel](../12.or-done-channel/README.md)

## Concept

Many data races disappear when concurrent code is designed around ownership
instead of unrestricted shared mutation. This lesson compares two forms.

### Partitioned memory

The result slice is allocated at its final length before workers start. Each
worker receives a pointer to a different element:

```go
processNum(num, &result[i])
```

The slice's backing array is shared, but the memory locations being written are
disjoint. No worker appends or changes the shared slice header. `wg.Wait`
finishes all writes before `main` reads the slice.

### Channel confinement

The balance actor is the only goroutine that may access `balance`. Other code
sends typed requests:

- withdraw an amount,
- deposit an amount, or
- query the current balance through a reply channel.

Because one goroutine processes the request stream serially, the state needs no
mutex.

## Request-response synchronization

The query request includes its own capacity-one reply channel. Receiving the
reply proves that the account processed the query and every request it received
before that query. The requester creates the channel, and the account sends
exactly once. Because the requester expects exactly one value, the reply
channel does not need a close signal.

The one-element buffer is a liveness boundary: once the account produces the
reply, it can continue even if that requester is delayed or abandons the
receive. An unbuffered reply would let one abandoned client wedge the sole
state owner and every later request. A reusable production API would also make
request and reply operations cancellation-aware.

The coordinator must first stop every client that can send a request. It then
closes the request channel and waits on `stopped`, making actor shutdown
observable rather than relying on process exit. Closing while another client
may still send would cause that client to panic.

## Ownership table

| Resource | Owner |
| --- | --- |
| `result[i]` | Worker assigned index `i` |
| Slice header and final read | `main` |
| Account balance | Account goroutine |
| Request channel closure | `BalanceMain` coordinator |
| Per-query reply channel | Requester creates a one-value buffer; account sends once |
| Stopped-channel closure | Account goroutine |

## Common mistake: check then act outside ownership

```go
if amount <= balance {
	go withdraw(amount)
}
```

Another operation can change `balance` after the check and before the
withdrawal. A mutex must cover the whole invariant, or the confined owner must
perform both the check and update as one request.

The example intentionally allows overdrafts because it demonstrates ownership,
not banking policy.

## Run it

```bash
go run ./11.confinement
go test ./11.confinement/...
```

## Experiments

- Add an account request that rejects overdrafts and returns an error.
- Start several clients that send deposits and verify the final balance.
- Replace the partitioned result with append and explain why a mutex is then
  required.

## References

- [Go memory model](https://go.dev/ref/mem)
- [Effective Go: sharing by communicating](https://go.dev/doc/effective_go#sharing)
- [Go race detector](https://go.dev/doc/articles/race_detector)
