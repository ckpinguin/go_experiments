package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var println = fmt.Println
var sprintf = fmt.Sprintf
var printf = fmt.Printf

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func fetcher(url string) <-chan string {
	out := make(chan string)
	fmt.Println("Running goroutine for", url, "with chan", out)
	go func() {
		defer wgRec.Done()
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		out <- "hey " + url
		close(out)
	}()
	return out
}

// See http://www.tapirgames.com/blog/golang-channel-closing:
// ...don't close a channel if the channel has multiple concurrent senders
// so we need to use a waitgroup for this, there is no solid way to
// use chans only
func fanIn(channels []<-chan string) <-chan string {
	var wgRec sync.WaitGroup
	broadcast := make(chan string)
	// The issue here is, that we do not know, which is the last sending
	// goroutine, so we cannot safely close the broadcast chan.
	go func() {
		for {
			for i := range channels {
				c := channels[i]
				// The first select here is to try to exit the goroutine
				// as early as possible. In fact, it is not essential
				// for this example, so it can be omitted.
				select {
				case s := <-c:
					broadcast <- s
				default:
				}
			}
		}
	}()
	return broadcast
}

func main() {
	urls := []string{"urlone", "urltwo", "urlthree"}
	channels := []<-chan string{}
	stop := make(chan struct{})
	totalTimeout := time.After(1 * time.Second)

	for _, url := range urls {
		fmt.Println("Starting goroutine for", url)
		ch := fetcher(url)
		channels = append(channels, ch)
	}
	r := fanIn(channels, stop)

	for {
		select {
		case s := <-r:
			println(s)
			// stop <- struct{}{}
		case <-totalTimeout: // signaling usage of a channel
			println("Timed out")
			stop <- struct{}{}
			close(stop)
		}
	}

	// wgRec.Wait()

}
