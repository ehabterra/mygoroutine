package main

import (
	"flag"
	"fmt"
)

func addAmount(amt int, ch chan int) {
	ch <- amt
}

func main() {
	balance := 0

	iterations := flag.Int("it", 10000000, "number of iterations")

	flag.Parse()

	fmt.Println("Iterations:", *iterations)
	amountChan := make(chan int)
	done := make(chan bool, 2)

	go func() {
		for {
			select {
			case amount := <-amountChan:
				balance += amount
			}
		}
	}()

	// increment
	go func() {
		fmt.Println("Increment")
		for i := 0; i < *iterations; i++ {
			addAmount(1, amountChan)
			// fmt.Println("Balance:", balance, i, "+")
			// runtime.Gosched()
		}
		done <- true
	}()

	// decrement
	go func() {
		fmt.Println("Decrement")
		for i := 0; i < *iterations; i++ {
			addAmount(-1, amountChan)
			// fmt.Println("Balance:", balance, i, "-")
			// runtime.Gosched()
		}
		done <- true
	}()

	<-done // increment
	<-done // decrement

	fmt.Printf("Final balance: %v\n", balance)

}
