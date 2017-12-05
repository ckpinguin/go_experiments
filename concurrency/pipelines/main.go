package main

import "fmt"

// returns a chan to read FROM
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

func main() {
	c := gen(2, 3, 4, 235335)
	out := sq(sq(c))
	for i := range out {
		fmt.Println(i)
	}
}
