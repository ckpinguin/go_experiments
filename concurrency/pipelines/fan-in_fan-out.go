package main

import "fmt"

// send given numbers into a channel
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
		close(out) // Important!
	}()
	return out
}

func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)

	for _, c := range chans {
		go func(c <-chan int) {
			out <- <-c
		}(c)
	}
	return out
}

func main() {
	in := gen(2, 3, 2e12, 35, -17)

	c1 := sq(in)
	c2 := sq(in)

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
