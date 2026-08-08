package main

import (
	"fmt"
	"sync"
)

var (
	mu      sync.RWMutex
	balance int
)

func Withdraw(amount int) bool {
	mu.Lock()
	defer mu.Unlock()

	deposit(-amount)
	if balance < 0 {
		deposit(amount)
		return false
	}

	return true
}

func Deposit(amount int) {
	mu.Lock()
	defer mu.Unlock()
	deposit(amount)
}

func Balance() int {
	mu.RLock()
	defer mu.RUnlock()
	return balance
}

func deposit(amout int) {
	balance += amout
}

func main() {
	go func() {
		Deposit(200)                // A1
		fmt.Println("=", Balance()) // A2
	}()

	go Deposit(100)

	for {
		fmt.Println("=", Balance()) // A2
		// time.Sleep(1 * time.Second)
	}
}

var (
	icons         map[string]string
	loadIconsOnce sync.Once
)

func loadIcons() {
	icons = map[string]string{
		"spades.png":   loadIcon("spades.png"),
		"hearts.png":   loadIcon("hearts.png"),
		"diamonds.png": loadIcon("diamonds.png"),
		"clubs.png":    loadIcon("clubs.png"),
	}
}

func loadIcon(url string) string {
	return url
}

func Icon(name string) string {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}
