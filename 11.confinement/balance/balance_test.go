package balance

import (
	"testing"
	"time"
)

func TestAccountSerializesRequestsAndStops(t *testing.T) {
	requests := make(chan accountRequest)
	stopped := make(chan struct{})
	go account(requests, stopped)

	sendRequest(t, requests, accountRequest{kind: withdraw, amount: 300})
	sendRequest(t, requests, accountRequest{kind: deposit, amount: 50})

	query, reply := newBalanceQuery()
	sendRequest(t, requests, query)

	if got, want := receiveBalance(t, reply), 750; got != want {
		t.Fatalf("balance = %d; want %d", got, want)
	}

	close(requests)
	waitForAccountStop(t, stopped)
}

func TestAbandonedReplyDoesNotBlockLaterRequestsOrShutdown(t *testing.T) {
	requests := make(chan accountRequest)
	stopped := make(chan struct{})
	go account(requests, stopped)

	abandonedQuery, _ := newBalanceQuery()
	sendRequest(t, requests, abandonedQuery)
	sendRequest(t, requests, accountRequest{kind: deposit, amount: 50})

	query, reply := newBalanceQuery()
	sendRequest(t, requests, query)
	if got, want := receiveBalance(t, reply), 1050; got != want {
		t.Fatalf("balance after abandoned reply = %d; want %d", got, want)
	}

	close(requests)
	waitForAccountStop(t, stopped)
}

func sendRequest(t *testing.T, requests chan<- accountRequest, request accountRequest) {
	t.Helper()

	select {
	case requests <- request:
	case <-time.After(time.Second):
		t.Fatal("account did not accept request")
	}
}

func receiveBalance(t *testing.T, reply <-chan int) int {
	t.Helper()

	select {
	case balance := <-reply:
		return balance
	case <-time.After(time.Second):
		t.Fatal("account did not reply to balance query")
		return 0
	}
}

func waitForAccountStop(t *testing.T, stopped <-chan struct{}) {
	t.Helper()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("account did not stop after requests closed")
	}
}
