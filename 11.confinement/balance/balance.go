package balance

import "fmt"

/*
var (
	balance = 1000
	mu      sync.Mutex
)

// approach 1
func withdraw(wg *sync.WaitGroup, amount int) {
	defer wg.Done()
	// if amount > balance {
	// 	return
	// }

	mu.Lock()
	balance -= amount
	mu.Unlock()
}
*/

//****************************************

type RequestType int

const (
	WITHDRAW RequestType = iota
	DEPOSITE
	BALANCE
)

type AccountReq struct {
	reqType RequestType
	amount  int
	c       chan int
}

func account(req <-chan AccountReq) {
	balance := 1000

	for msg := range req {
		switch msg.reqType {
		case WITHDRAW:
			balance -= msg.amount
		case DEPOSITE:
			balance += msg.amount
		case BALANCE:
			msg.c <- balance
		}
	}
}

func BalanceMain() {
	accChannel := make(chan AccountReq)
	defer close(accChannel)

	go account(accChannel)

	accChannel <- AccountReq{
		reqType: WITHDRAW,
		amount:  300,
	}

	accChannel <- AccountReq{
		reqType: DEPOSITE,
		amount:  50,
	}

	balanceChannel := make(chan int)

	accChannel <- AccountReq{
		reqType: BALANCE,
		c:       balanceChannel,
	}

	fmt.Println("Balance: ", <-balanceChannel)
}
