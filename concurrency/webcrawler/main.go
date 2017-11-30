// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// I deliberately copied and learned from the official solution (copyright above)

package main

import (
	"errors"
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

// fetched tracks URLs that have been (or are being) fetched.
// The lock must be held while reading from or writing to the map.
// See http://golang.org/ref/spec#Struct_types section on embedded types.type Cache struct {
var fetched = struct {
	m map[string]error
	sync.RWMutex
}{m: make(map[string]error)}

// var fetchedSync = sync.Map[string]error{}

var errLoading = errors.New("url load in progress")

func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		fmt.Printf("<- Done with %v, reached depth 0.\n", url)
		return
	}
	// Start reading
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

	// Fetch stuff
	body, urls, err := fetcher.Fetch(url)

	// Update status
	fetched.Lock()
	fetched.m[url] = err // not loading anymore (nil or a "real" error)
	fetched.Unlock()

	if err != nil {
		fmt.Printf("<- Error on %v: %v\n", url, err)
		return
	}

	fmt.Printf("Found: %s %q\n", url, body)

	// every Crawl goroutine call has it's own done chan
	done := make(chan bool)

	// start the goroutines for every sub-url found under the current url
	for i, u := range urls {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i+1, len(urls), url, u)
		go func(url string) {
			// Recursion here
			Crawl(url, depth-1, fetcher)
			done <- true // signal goroutine is done
		}(u) // <= url
	}
	// synchronize/wait for the goroutines to finish
	for i, u := range urls {
		fmt.Printf("<- [%v] %v/%v waiting for child %v.\n", url, i+1, len(urls), u)
		<-done // block until all children are done
	}
	fmt.Printf("<- Done with %v\n", url)
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)

	fmt.Println("Fetching stats\n----------------")
	for url, err := range fetched.m {
		if err != nil {
			fmt.Printf("%v failed: %v\n", url, err)
		} else {
			fmt.Printf("%v was fetched\n", url)
		}
	}
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	// comma, ok idiom to identify missing values (discriminate them from zero-values)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
// This could as well serve as mock data
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
