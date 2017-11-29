package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	// Simple and elegant from https://golang.org/doc/play/tree.go
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

func Walker(t *tree.Tree) chan int {
	ch := make(chan int)
	// This must be an anon. self exec. func, because otherwise the channel
	// would be closed before it is returned and thus not giving any values!
	go func() {
		Walk(t, ch)
		close(ch)
	}()
	return ch
}

// Same determines whether the trees
// t1 and t2 contain the same values.
// It reads values from two Walkers
// that run simultaneously, and returns true
// if t1 and t2 have the same contents.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := Walker(t1)
	ch2 := Walker(t2)
	for {
		// comma, ok idiom to discriminate a 0 (no nil for int!) from a missing value
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if !ok1 || !ok2 {
			// the only case in which a "true" can be returned, is
			// when both node-value fields are empty (ok set to false for both)
			return ok1 == ok2
		}
		if v1 != v2 {
			// in any case where there is a difference in values, the func
			// breaks and returns false
			break
		}
	}
	return false
}

func main() {
	t1 := tree.New(1)
	t2 := tree.New(1)
	t3 := tree.New(2)
	ch := Walker(t3)
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println(Same(t1, t2))
}
