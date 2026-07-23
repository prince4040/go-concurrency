# Lesson 6: Done Channel

Previous: [For/select](../5.for-select/README.md) ·
[Course index](../README.md) · Next: [Pipeline](../7.pipeline/README.md)

## Concept

Closing a channel broadcasts an event to every receiver. A cancellation-only
channel conventionally uses `chan struct{}` because no value is needed:

```go
done := make(chan struct{})
close(done)
```

Every `<-done` can then proceed immediately.

## Walkthrough

`fibonacci` owns a number stream and produces until either:

- a consumer is ready for the next value, or
- the caller closes `done`.

After reading ten numbers, `main` closes `done` and ranges over the output until
the producer closes it. That final range is an explicit wait for cleanup.

## Ownership

| Channel | Sender/closer | Receiver |
| --- | --- | --- |
| `done` | `main` closes; nobody sends | Fibonacci producer |
| `numbers` | Fibonacci producer sends and closes | `main` |

Only the coordinator closes `done`, and it does so once. Sending a value on
`done` would wake only one receiver; closing wakes all current and future
receivers.

## Common mistake: busy cancellation polling

```go
for {
	select {
	case <-done:
		return
	default:
		fmt.Println("working")
	}
}
```

Without real blocking work or a deliberate wait, this loop spins as fast as
possible. Put the work's blocking operation in the same select, or introduce a
timer/ticker appropriate to the task.

## Cancellation is cooperative

Closing `done` does not forcibly terminate a goroutine. The goroutine must
observe the signal at every operation where it might otherwise block
indefinitely.

Cancellation also has no priority. If `done` is closed at the same instant a
consumer is ready for `numbers`, both select cases are ready and either may
run. The producer can therefore deliver one final value that raced with the
cancellation request. Consumers must treat cancellation as a request to stop,
not as a promise that no more values can arrive.

## Run it

```bash
go run ./6.done-channel
```

## Experiments

- Start two Fibonacci producers with the same `done` channel.
- Stop after one value and confirm both producer outputs close.
- Remove the send-side cancellation case and explain the possible leak.

## References

- [Go memory model: channel closing](https://go.dev/ref/mem)
- [Go concurrency patterns: pipelines and cancellation](https://go.dev/blog/pipelines)
