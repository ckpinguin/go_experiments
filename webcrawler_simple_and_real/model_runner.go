package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func fetcher(url string) <-chan string {
	out := make(chan string)
	fmt.Println("Running goroutine for", url, "with chan", out)
	go func() {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		out <- "hey " + url
		close(out)
	}()
	return out
}

// See http://www.tapirgames.com/blog/golang-channel-closing:
// ...don't close a channel if the channel has multiple concurrent senders
// so we need to use a waitgroup for this, there is no solid way to
// use chans only
func fanIn(cs []<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		for s := range c {
			out <- s
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	urls := []string{"urlone", "urltwo", "urlthree"}
	channels := []<-chan string{}

	for _, url := range urls {
		fmt.Println("Starting goroutine for", url)
		ch := fetcher(url)
		channels = append(channels, ch)
	}
	r := fanIn(channels)

	totalTimeout := time.After(5 * time.Second)
loop:
	for {
		select {
		case s, ok := <-r:
			if !ok {
				break loop
			}
			fmt.Println(s)
		case <-totalTimeout: // signaling usage of a channel
			fmt.Println("Timed out")
			break loop
		}
	}
}
