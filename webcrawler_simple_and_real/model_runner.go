package main

import (
	"fmt"
	"math/rand"
	"time"
)

var println = fmt.Println
var sprintf = fmt.Sprintf
var printf = fmt.Printf

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func fetcher(url string) <-chan string {
	ch := make(chan string)
	fmt.Println("Running goroutine for", url, "with chan", ch)
	go func() {
		ch <- "hey " + url
	}()
	return ch
}

func fanIn(channels []<-chan string) <-chan string {
	broadcast := make(chan string)
	go func() {
		for {
			for i := range channels {
				c := channels[i]
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
	totalTimeout := time.After(2 * time.Second)

	for _, url := range urls {
		fmt.Println("Starting goroutine for", url)
		channels = append(channels, fetcher(url))
	}
	r := fanIn(channels)
	for {
		select {
		case s := <-r:
			println(s)
		case <-totalTimeout:
			println("Timed out")
			return
		}
	}
}
