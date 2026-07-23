package balance

import "fmt"

/*
The account goroutine owns balance for its entire lifetime. Requests arrive on
one channel, so withdrawals, deposits, and reads are processed serially without
a mutex. A reply channel lets a balance query synchronize with all requests
that were sent before it.
*/

type requestType int

const (
	withdraw requestType = iota
	deposit
	getBalance
)

type accountRequest struct {
	kind   requestType
	amount int
	reply  chan<- int
}

func newBalanceQuery() (accountRequest, <-chan int) {
	reply := make(chan int, 1)
	return accountRequest{kind: getBalance, reply: reply}, reply
}

func account(requests <-chan accountRequest, stopped chan<- struct{}) {
	defer close(stopped)

	balance := 1000

	for request := range requests {
		switch request.kind {
		case withdraw:
			balance -= request.amount
		case deposit:
			balance += request.amount
		case getBalance:
			request.reply <- balance
		default:
			panic("unknown account request")
		}
	}
}

func BalanceMain() {
	requests := make(chan accountRequest)
	stopped := make(chan struct{})

	go account(requests, stopped)

	requests <- accountRequest{
		kind:   withdraw,
		amount: 300,
	}

	requests <- accountRequest{
		kind:   deposit,
		amount: 50,
	}

	query, reply := newBalanceQuery()
	requests <- query

	fmt.Println("Balance:", <-reply)

	close(requests)
	<-stopped
}
