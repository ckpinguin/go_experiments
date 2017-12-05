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

func Crawler(urls []string) chan string {
	cs := make([]chan string, len(urls))
	for i, url := range urls {
		cs[i] = make(chan string)
		log.Println("Starting runner for", url)
		go crawl(url, cs[i])
	}
	s := fanIn(cs...)
	return s
}

func fanIn(inputs ...chan string) chan string {
	c := make(chan string)
	for i := range inputs {
		input := inputs[i] // New instance of 'input' for each loop.
		go func() {
			for {
				c <- <-input
			}
		}()
	}
	return c
}

func crawl(url string, c chan string) {
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
		fmt.Println("Starting parse...")
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					hasProto := strings.Index(a.Val, "http") == 0
					if hasProto {
						c <- a.Val
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

	// z := html.NewTokenizer(resp.Body)

	// for {
	// 	tt := z.Next()
	// 	var t html.Token
	// 	switch tt {
	// 	case html.ErrorToken:
	// 		break
	// 	case html.StartTagToken:
	// 		t = z.Token()
	// 		isAnchor := t.Data == "a"
	// 		if !isAnchor {
	// 			continue
	// 		}
	// 		link, err := getHref(t)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		hasProto := strings.Index(link, "http") == 0
	// 		if hasProto {
	// 			// log.Printf("Found on %v: %v\n", url, link)
	// 			fetched.Lock()
	// 			fetched.m[url] = nil // not loading anymore (nil or a "real" error)
	// 			fetched.m[link] = nil
	// 			fetched.Unlock()
	// 			c <- link
	// 			break
	// 		}
	// 	}
	// }
	close(c)
	return
}

func main() {
	urls := os.Args[1:]

	c := Crawler(urls)
	timeout := time.After(2 * time.Second)

	for {
		select {
		case s := <-c:
			fmt.Println(" - ", s)
		case <-timeout:
			fmt.Println("Timeout!")
			return
		}
	}

	fmt.Println("Found", len(fetched.m), "unique urls:")
	for url := range fetched.m {
		fmt.Println(" - ", url)
	}

}

// Extract href string from a Token
func getHref(t html.Token) (href string, err error) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			err = nil
			break
		}
		err = errors.New("No href attribute found")
	}
	return
}
