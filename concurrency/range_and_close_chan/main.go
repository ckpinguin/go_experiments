package main

import "fmt"

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		// c <- x
		x, y = y, x+y
	}
	c <- x // just the result please
	close(c)
}

func main() {
	c := make(chan int, 92)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}
