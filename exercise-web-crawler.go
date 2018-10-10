package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Cache stores the fetched urls
type Cache struct {
    url map[string]bool
	mux sync.Mutex
}

// Set assigns a key concurrently
func (c *Cache) Set(key string) {
    c.mux.Lock()
	defer c.mux.Unlock()
	c.url[key] = true
}

// Get checks if key is present
func (c *Cache) Get(key string) (bool) {
    c.mux.Lock()
	defer c.mux.Unlock()
    _, ok := c.url[key]
	return ok
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, fetcher fakeFetcher, cache *Cache, ch chan string) {
    // Check if the url is already in the cache
    if cache.Get(url) {
	    return
	}

	cache.Set(url)
	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		ch <- fmt.Sprintf("%v\n", err)
		return
	}

	ch <- fmt.Sprintf("found: %s %q\n", url, body)

	for _, u := range urls {
		go Crawl(u, fetcher, cache, ch)
	}

	return
}

func main() {
    cache := Cache{url: make(map[string]bool)}
    ch := make(chan string)
    defer close(ch)
	go Crawl("https://golang.org/", fetcher, &cache, ch)

	for i := 0; i < len(fetcher); i++ {
	    fmt.Print(<-ch)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}

	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

