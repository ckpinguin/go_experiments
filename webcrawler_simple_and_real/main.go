package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var fetched = struct {
	m map[string]error
	sync.RWMutex
}{m: make(map[string]error)}

var errLoading = errors.New("url load in progress")

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
			break
		}
	}
	return
}

func crawl(url string, done chan bool) {
	// every Crawl goroutine call has it's own done chan
	r, err := http.Get(url)
	if err != nil {
		log.Fatalln("ERROR: Could not read from", url, err)
	}
	defer r.Body.Close()

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

	z := html.NewTokenizer(r.Body)
	for {
		tt := z.Next()
		var t html.Token
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t = z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			ok, link := getHref(t)
			if !ok {
				continue
			}
			hasProto := strings.Index(link, "http") == 0
			if hasProto {
				// log.Printf("Found on %v: %v\n", url, link)
				fetched.Lock()
				fetched.m[url] = nil // not loading anymore (nil or a "real" error)
				fetched.m[link] = nil
				fetched.Unlock()
				// o <- url
			}
		}
	}
	done <- true
}

func main() {
	// const url = "http://www.apodro.ch"
	//foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	// res := make(chan string)
	done := make(chan bool)
	for _, url := range seedUrls {
		log.Println("Starting runner for", url)
		go crawl(url, done)
	}
	for _, url := range seedUrls {
		log.Printf("Waiting for child on %v\n", url)
		<-done
		log.Printf("Child on %v has finished\n", url)
	}

	fmt.Println("Found", len(fetched.m), "unique urls:")
	for url := range fetched.m {
		fmt.Println(" - ", url)
	}

}
