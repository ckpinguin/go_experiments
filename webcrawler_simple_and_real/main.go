package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// bookkeeping of loaded urls
var fetched = struct {
	m map[string]error
	sync.RWMutex
}{m: make(map[string]error)}

var errLoading = errors.New("url load in progress")

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

func fetcher(url string) <-chan string {
	out := make(chan string)
	// do this async
	go func() {
		fetched.RLock()
		if _, ok := fetched.m[url]; ok {
			fetched.RUnlock()
			fmt.Printf("<- Done with %v, already fetched.\n", url)
			return
		}
		fetched.RUnlock()

		// Start writing
		// Mark for loading
		fetched.Lock()
		fetched.m[url] = errLoading
		fetched.Unlock()
		// End of "transaction"
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln("ERROR: Could not read from", url, err)
		}
		defer resp.Body.Close()
		var s []byte
		resp.Body.Read(s)

		doc, err := html.Parse(resp.Body)
		var f func(*html.Node)
		f = func(n *html.Node) {
			// fmt.Println("Starting parse...")
			if n.Type == html.ElementNode && n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key == "href" {
						fmt.Println(a.Val)
						hasProto := strings.Index(a.Val, "http") == 0
						if hasProto {
							out <- a.Val
						}
						break // break after the first (hopefully only) href
					}
				}
			}
			// Recurse
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
		close(out)
	}()
	return out
}

func main() {
	urls := os.Args[1:]
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

	fmt.Println("Found", len(fetched.m), "unique urls:")
	for url := range fetched.m {
		fmt.Println(" - ", url)
	}

}
