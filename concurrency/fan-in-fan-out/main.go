package main

import (
	"fmt"
	"sync"
)

func main() {
	in := gen(2, 3, 4, 5, 34535)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)
	c3 := sq(in)
	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2, c3) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}

// Merging multiple chan's into one
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Copy the values from c to out until c is closed,
	// then remove from itself WaitGroup using Done()
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs)) // for every chan of cs

	// Start an output goroutine for each input chan in cs.
	for _, c := range cs {
		go output(c)
	}

	// Close out chan if all input chan's have been handled
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}
