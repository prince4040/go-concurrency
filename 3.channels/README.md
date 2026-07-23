# Lesson 3: Channels

Previous: [WaitGroup](../2.waitGroup/README.md) ·
[Course index](../README.md) · Next: [Select](../4.select/README.md)

## Concept

A channel transfers typed values between goroutines. It is both communication
and synchronization.

```go
c := make(chan string) // unbuffered
c <- "data"             // send
value := <-c            // receive
```

On an unbuffered channel, the send waits for a receiver and the receive waits
for a sender. Their handshake creates a synchronization point.

## Walkthrough

1. `main` creates an unbuffered channel.
2. The producer goroutine attempts to send `"data"`.
3. The send and `main`'s receive meet; `main` gets the value with `open=true`.
4. The producer closes its output.
5. A second receive returns the zero value and `open=false`.

Closing does not erase buffered or already-sent values. Receivers consume those
values before observing `open=false`.

## Ownership and direction

The parameter `chan<- string` documents that `send` may only send. A
`<-chan string` parameter would be receive-only.

The producer closes the channel because it knows when the final send is done.
Receivers generally must not close a channel they do not own.

## Common mistake: receiver-side close

```go
value := <-c
close(c) // unsafe if a producer might send again
```

Closing while another goroutine sends can panic. A channel does not need to be
closed after every use; close it only to communicate that no more values will
ever be sent.

## Buffered channels

`make(chan T, n)` allows up to `n` values to wait in the channel. A send blocks
only when the buffer is full. A buffer changes when goroutines synchronize; it
does not remove the need for ownership or cancellation.

## Run it

```bash
go run ./3.channels
```

## Experiments

- Change the channel to `make(chan string, 1)`.
- Send two values and receive until the channel closes.
- Remove the goroutine and observe the unbuffered-send deadlock.

## References

- [Go specification: channel types](https://go.dev/ref/spec#Channel_types)
- [Go memory model: channel communication](https://go.dev/ref/mem)
- [Effective Go: channels](https://go.dev/doc/effective_go#channels)
