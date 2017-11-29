package main

import "testing"

func BenchmarkCrawl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := make(chan string)
		go Crawl("http://golang.org/", 4, fetcher, result)
	}
}
func BenchmarkCrawlMerged(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := make(chan string)
		go CrawlMerged("http://golang.org/", 4, fetcher, result)
	}
}

func TestCrawl(t *testing.T) {
	type args struct {
		url     string
		depth   int
		fetcher Fetcher
		ret     chan string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Crawl(tt.args.url, tt.args.depth, tt.args.fetcher, tt.args.ret)
		})
	}
}
