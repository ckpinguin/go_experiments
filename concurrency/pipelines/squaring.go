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

func main() {
	c := gen(1, 2, 3, -2, 4e8, 5)
	out := sq(sq(c))

	// consume out values
	for v := range out {
		fmt.Println(v)
	}
	// Set up the pipeline directly and consume the output.
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n) // 16 then 81
	}
}
