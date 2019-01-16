package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, wg sync.WaitGroup) {
	// X TODO: Fetch URLs in parallel.
	// X TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	wg.Add(1)
	defer wg.Done()
	time.Sleep(100* time.Millisecond)
	
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, wg)
	}
	return
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

// cachedFetcher is a parallel-safe fetcher with a cache
type CachedFetcher struct {
	fetcher Fetcher
	cache   map[string]CachedResult
	lock    sync.RWMutex
}

type CachedResult struct {
	err error
	fakeResult
}

func (cf CachedFetcher) Fetch(url string) (body string, urls []string, err error) {
	cf.lock.RLock()
	cachedResult, ok := cf.cache[url]
	cf.lock.RUnlock()
	if ok {
		body, urls, err = cachedResult.body, cachedResult.urls, cachedResult.err
		fmt.Printf("Cache Hit: ")
		return
	}
	cf.lock.Lock()
	defer cf.lock.Unlock()
	body, urls, err = cf.fetcher.Fetch(url)
	cachedResult = CachedResult{err: err, fakeResult: fakeResult{body, urls}}
	fmt.Printf("-- created cache for %v = %v\n", url, cachedResult)
	fmt.Printf("Cache NEW: ")
	cf.cache[url] = cachedResult
	return
}

var cachedFetcher = CachedFetcher{
	fetcher: fetcher,
	cache:   make(map[string]CachedResult),
	lock: sync.RWMutex{},
}

func main() {
	wg := sync.WaitGroup{}
	
	Crawl("https://golang.org/", 4, cachedFetcher, wg)
	time.Sleep(time.Second)  // TODO this wait is really annoying and has to be avoided
	// also note that we have no way to communicate the crawler results to the rest of the world, this only works
	// because fmt.Printf is basically a queue for human eyes :-) 
	wg.Wait()
	
	fmt.Printf("Created %v cache entries", len(cachedFetcher.cache))
}
