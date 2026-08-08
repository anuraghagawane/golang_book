package main

import (
	"fmt"
	"time"
)

var (
	deposits = make(chan int)
	balances = make(chan int)
)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func main() {
	go teller()

	go func() {
		Deposit(200)                // A1
		fmt.Println("=", Balance()) // A2
	}()

	go Deposit(100)

	time.Sleep(5 * time.Second)
}
