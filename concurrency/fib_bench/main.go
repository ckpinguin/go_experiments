package main

import (
	"fmt"
	"time"
)

func fib(n int) uint64 {
	if n <= 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}

func fibIt(n int) uint64 {
	// defer timeTrack(time.Now(), "Iterative Fibonacci of "+strconv.Itoa(n))
	var x, y uint64
	x, y = 0, 1
	for i := 0; i < n; i++ {
		// c <- x
		x, y = y, x+y
	}
	return x
}

func fibP(n int) uint64 {
	// defer timeTrack(time.Now(), "Peter's Iterative Fibonacci of "+strconv.Itoa(n))
	f := uint64(0)
	a, b := uint64(0), uint64(1)
	for i := 0; i <= n; i++ {
		f, a, b = a, b, a+b
		if a > b {
			break
		}
	}
	return f
}

func main() {
	var start time.Time
	var elapsed time.Duration

	start = time.Now()
	fibP(46)
	elapsed = time.Since(start)
	fmt.Println("Peter's Iterative Fibonacci of 46 took", elapsed)

	start = time.Now()
	fibIt(46)
	elapsed = time.Since(start)
	fmt.Println("Iterative Fibonacci of 46 took", elapsed)

	start = time.Now()
	fib(46)
	elapsed = time.Since(start)
	fmt.Println("Recursive Fibonacci of 46 took", elapsed)
}
