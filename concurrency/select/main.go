package main

import "fmt"

func fibonacci(c chan uint64, quit chan uint64) {
	var x, y uint64
	x, y = 0, 1
	for {
		select {
		case c <- x: // send x to chan c (if it is ready to take a value)
			x, y = y, x+y
		case <-quit: // receive a quit signal (read from input chan)
			fmt.Println("quit")
			return
		}
	}
}

// Mnemonic: chan<- means "something to SEND TO the channel"
// <-chan means "something to RECEIVE FROM a channel"
func main() {
	c := make(chan uint64)    // input chan (read FROM the chan)
	quit := make(chan uint64) // output chan (write TO the chan)
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println(<-c) // read from input chan
		}
		quit <- 0 // signal quit
	}()
	fibonacci(c, quit)
}
