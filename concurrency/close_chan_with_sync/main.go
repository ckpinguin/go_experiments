package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, more := <-jobs
			if more {
				io.WriteString(os.Stdout, "received job "+strconv.Itoa(j)+"\n")
				// fmt.Println("BUF received job", j)
			} else {
				fmt.Println("all jobs received")
				done <- true
				return
			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		io.WriteString(os.Stdout, "sent job "+strconv.Itoa(j)+"\n")
		// fmt.Println("BUF sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	<-done

}
