package main

import "fmt"

func main() {
	// naturals := make(chan int)
	// squares := make(chan int)
	//
	// go func() {
	// 	for x := 0; x <= 100; x++ {
	// 		naturals <- x
	// 	}
	// 	close(naturals)
	// }()
	//
	// go func() {
	// 	for x := range naturals {
	// 		squares <- x * x
	// 	}
	// 	close(squares)
	// }()
	//
	// for sq := range squares {
	// 	fmt.Println(sq)
	// }

	naturals := make(chan int)
	squares := make(chan int)
	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

func counter(out chan<- int) {
	for x := 0; x <= 100; x++ {
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for v := range in {
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}
