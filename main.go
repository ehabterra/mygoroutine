package main

import (
	"flag"
	"fmt"
	"runtime"
)

func main() {
	balance := 0

	finishedIncremental := false
	finishedDecremental := false

	iterations := flag.Int("it", 10000000, "number of iterations")

	flag.Parse()

	fmt.Println("Iterations:", *iterations)

	amountChan := make(chan int)

	defer close(amountChan)

	incrementDone := make(chan bool)

	// increment
	go func() {
		fmt.Println("Increment")
		for i := 0; i < *iterations; i++ {
			amountChan <- 1
			// fmt.Println("Balance:", balance, i, "+")
			runtime.Gosched()
		}
		incrementDone <- true
	}()

	decrementDone := make(chan bool)

	// decrement
	go func() {
		fmt.Println("Decrement")
		for i := 0; i < *iterations; i++ {
			amountChan <- -1
			// fmt.Println("Balance:", balance, i, "-")
			runtime.Gosched()
		}
		decrementDone <- true
	}()

	// print final balance and close channels
	cleanUp := func() {
		close(incrementDone)
		close(decrementDone)
		fmt.Printf("Final balance: %v\n", balance)
	}

	for {
		select {
		case amount := <-amountChan:
			balance += amount
		case finishedIncremental = <-incrementDone:
			if finishedIncremental && finishedDecremental {
				cleanUp()
				return
			}
		case finishedDecremental = <-decrementDone:
			if finishedIncremental && finishedDecremental {
				cleanUp()
				return
			}
		}
	}

}
