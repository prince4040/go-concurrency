# Learn Go Concurrency by Running It

This repository is a progressive, executable course on Go concurrency. Each
numbered directory contains one small program, an explanation of the mental
model behind it, common failure modes, and experiments you can try.

The examples favor correctness and explicit ownership over cleverness. Every
canonical program terminates cleanly; intentionally broken approaches appear
only as labeled snippets in the lesson notes.

## Prerequisites

- Basic Go syntax: functions, loops, slices, interfaces, pointers, closures,
  and `defer`
- Basic type-parameter syntax such as `[T any]`; Lesson 8 includes a short
  concurrency-focused refresher
- Go 1.26.4 or newer, matching [`go.mod`](go.mod)
- A terminal from the repository root

No third-party dependencies are required.

## Learning path

| Lesson | Topic | Central question |
| --- | --- | --- |
| [1](1.routine/README.md) | Goroutines | How does work run independently? |
| [2](2.waitGroup/README.md) | WaitGroup | How does a caller wait for a set of tasks? |
| [3](3.channels/README.md) | Channels | How do goroutines communicate and synchronize? |
| [4](4.select/README.md) | Select | How do we wait on several channel operations? |
| [5](5.for-select/README.md) | For/select | How do we multiplex streams until all are closed? |
| [6](6.done-channel/README.md) | Done channel | How do we broadcast cancellation? |
| [7](7.pipeline/README.md) | Pipeline | How do stages compose without leaking goroutines? |
| [8](8.generator/README.md) | Generator | How do we build and bound an on-demand stream? |
| [9](9.fan-in-out/README.md) | Fan-out/fan-in | How do workers distribute and merge work? |
| [10](10.mutex/README.md) | Mutex | How do we protect shared mutable state? |
| [11](11.confinement/README.md) | Confinement | How can ownership remove the need for a lock? |
| [12](12.or-done-channel/README.md) | Or-done channel | How can a range loop remain cancellable? |
| [13](13.context/README.md) | Context | How do cancellation and deadlines follow a call tree? |

Follow the lessons in order on a first pass. Lessons 1–6 introduce the
primitives; lessons 7–13 compose them into reusable patterns.

## Running the examples

Run one lesson from the repository root:

```bash
go run ./1.routine
go run ./7.pipeline
go run ./13.context
```

Run the repository checks:

```bash
go test ./...
go vet ./...
go test -race ./...
```

The race detector observes only code paths that execute. For a runnable lesson,
you can also execute its `main` function under the detector:

```bash
go run -race ./10.mutex
go run -race ./11.confinement
```

Every numbered lesson has at least one behavioral test. The tests follow the
same rules learners should use in concurrent programs:

- Assert exact values or membership, but not incidental goroutine ordering.
- Assert order only when one producer or pipeline stage guarantees it.
- Bound lifecycle waits by one second so a leak fails promptly.
- Use channels, locks, or WaitGroups for synchronization—never `time.Sleep`.

Passing the race detector means no race was observed on the executed paths; it
is not proof that every possible path is race-free.

## The ownership model

Most channel-based examples follow the same lifecycle:

1. A function creates an output channel.
2. That function starts the only goroutine allowed to send on it.
3. The sending goroutine closes the output after its final send.
4. Consumers receive until closure, but never close that channel.
5. The coordinator creates and closes a cancellation signal.

This gives every channel one clear closer. Sending on a closed channel or
closing a channel twice panics, so ambiguous ownership is a correctness bug.

For shared memory, identify the synchronization boundary explicitly:

- A mutex protects every concurrent access to the same mutable state.
- Partitioning gives workers disjoint memory locations.
- Confinement gives one goroutine exclusive access and communicates requests.

## What Go guarantees

- Starting a goroutine does not wait for it to finish.
- A channel send is synchronized before its matching receive completes.
- Closing a channel is synchronized before a receive that observes the closure.
- `WaitGroup.Wait` returns only after its tracked tasks finish.
- `Mutex.Unlock` synchronizes before a later successful lock of the same mutex.

Go does not guarantee:

- Which runnable goroutine executes first
- The order in which concurrent results arrive
- That adding workers improves performance
- That a `select` prefers the first case in source order

Tests should therefore assert values, ownership, and completion properties—not
an accidental print order.

## Glossary

**Concurrency:** Structuring a program as independently progressing tasks.

**Parallelism:** Executing work at the same instant on multiple processing
resources. Concurrent code may or may not run in parallel.

**Blocking:** Waiting until an operation can proceed, such as an unbuffered send
waiting for a receiver.

**Data race:** Concurrent access to one memory location where at least one
access is a write and the accesses are not synchronized.

**Goroutine leak:** A goroutine that remains blocked with no possible useful
way to finish.

**Backpressure:** A slow consumer causing upstream senders to wait instead of
allowing work to grow without a bound.

**Fan-out:** Multiple workers receive from one input and divide its values.

**Fan-in:** Multiple inputs are merged into one output.

## How to study each lesson

1. Read its README and predict the output properties.
2. Run the example several times.
3. Make one suggested experiment.
4. Explain who owns and closes each channel.
5. Run the example with `-race` after changing shared-state code.

Avoid using `time.Sleep` to make a concurrency bug “go away.” A delay changes
probability, not synchronization.

These lessons require Go 1.26.4, whose `for range` semantics create new
iteration variables on every iteration. When reading advice written for Go
versions before 1.22, remember that closures no longer share an ordinary range
iteration variable. Variables declared outside the loop can still be shared
and require synchronization.

## Official references

- [The Go memory model](https://go.dev/ref/mem)
- [The Go language specification: statements and select](https://go.dev/ref/spec)
- [Effective Go: concurrency and channels](https://go.dev/doc/effective_go#concurrency)
- [`sync` package documentation](https://pkg.go.dev/sync)
- [Go concurrency patterns: pipelines and cancellation](https://go.dev/blog/pipelines)
- [`context` package documentation](https://pkg.go.dev/context)
- [Go data race detector](https://go.dev/doc/articles/race_detector)
