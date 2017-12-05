package main

import (
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	id int
}

func init() {
	s := time.Now().UTC().UnixNano()
	rand.Seed(s)
}

const MaxOutstanding = 10

var sem = make(chan int, MaxOutstanding)

func Serve(queue <-chan *Request) {
	for req := range queue {
		sem <- 1
		go func(req *Request) {
			process(req)
			<-sem
		}(req)
	}
}

// Funcs without return value are good for goroutines or inside goroutines
func process(r *Request) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	io.Copy(os.Stdout, strings.NewReader("process: "+strconv.Itoa(r.id)+"\n"))
}

func main() {

	clientRequests := make(chan *Request)
	//  quit := make(chan bool)
	go Serve(clientRequests)

	for i := 0; i < 100; i++ {
		// fmt.Println("doing", i)
		clientRequests <- &Request{i}
	}

	time.Sleep(time.Second * 6)
	// quit <- true
}
