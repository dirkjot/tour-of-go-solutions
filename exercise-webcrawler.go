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
// 
// - string: page to visit
// - depth:  recursion depth, decreasing, will be capped at 0
// - fetcher:  fetcher interface for (fake) retrieval of pages
// - control:  recursionController interface to communicate results and progress
//
func Crawl(url string, depth int, fetcher Fetcher, control recursionController) {
	defer func(control recursionController) { 
		control.Done(1)
	}(control)
	
	time.Sleep(100 * time.Millisecond) // mimic internet delay
	if depth <= 0 {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		control.results <- err.Error()
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	control.results <- fmt.Sprintf("Found: '%s' with urls: %v", body, urls)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher, control)
	}
	control.Add(len(urls))
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
	return "", nil, fmt.Errorf("Not found: %s", url)
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
// it contains a wrapped fetcher, for filling the cache.  Internally, 
// it uses a cache (map) and a RWMutex to allow multiple readers but only one writer
//
type CachedFetcher struct {
	fetcher Fetcher
	cache   map[string]CachedResult
	lock    *sync.RWMutex
}

// CachedResult data structure for cachedFetcher, wraps fakeResult 
// using go struct embedding https://golang.org/doc/effective_go.html#embedding
//
type CachedResult struct {
	err error
	fakeResult
}

// Fetch for CachedFetcher, returns a cached value, if missing, calls wrapped fetcher. 
// 
func (cf CachedFetcher) Fetch(url string) (body string, urls []string, err error) {
	cf.lock.RLock()
	cachedResult, ok := cf.cache[url]
	cf.lock.RUnlock()
	if ok {
		body, urls, err = cachedResult.body, cachedResult.urls, cachedResult.err
		return
	}
	cf.lock.Lock()
	defer cf.lock.Unlock()
	body, urls, err = cf.fetcher.Fetch(url)
	cachedResult = CachedResult{err: err, fakeResult: fakeResult{body, urls}}
	// fmt.Printf("-- created cache for %v = %v\n", url, cachedResult)
	cf.cache[url] = cachedResult
	return
}

// cachedFetcher : initialize our one and only CachedFetcher
var cachedFetcher = CachedFetcher{
	fetcher: fetcher,
	cache:   make(map[string]CachedResult),
	lock:    &sync.RWMutex{},
}

// recursionController is a data structure with three channels to control our Crawl recursion.
// Tried to use sync.waitGroup in a previous version, but I was unhappy with the mandatory sleep.
// The idea is to have three channels, counting the outstanding calls (children), completed calls 
// (done) and results (results).  Once outstanding calls == completed calls we are done (if you are
// sufficiently careful to signal any new children before closing your current one, as you may be the last one).
//
type recursionController struct {
	results  chan string
	children chan int
	done     chan int
}

// instead of instantiating one instance, as we did above, use a more idiomatic Go solution
func NewRecursionController() recursionController {
	// we buffer results to 1000, so we cannot crawl more pages than that.  
	return recursionController{make(chan string, 1000), make(chan int), make(chan int)}
}

// recursionController.Add: convenience function to add children to controller (similar to waitGroup)
func (rc recursionController) Add(children int) {
	rc.children <- children
}

// recursionController.Done: convenience function to remove children from controller (similar to waitGroup)
func (rc recursionController) Done(children int) {
	rc.done <- children
}

// recursionController.Wait will wait until all children are done
func (rc recursionController) Wait() {
	fmt.Println("Controller waiting...")
	var children, done int
	for {
		select {
		case childrenDelta := <-rc.children:
			children += childrenDelta
			// fmt.Printf("children found %v total %v\n", childrenDelta, children)
		case <-rc.done:
			done += 1
			// fmt.Println("done found", done)
		default:
			if done > 0 && children == done {
				fmt.Printf("Controller exiting, done = %v, children =  %v\n", done, children)
				close(rc.results)
				return
			}
		}
	}
}

func main() {
	control := NewRecursionController()
	go control.Add(1) // tell controller we are going to launch our first child 
	go Crawl("https://golang.org/", 4, cachedFetcher, control)
	
	// let the go routines flow
	control.Wait()

	fmt.Printf("Created %v cache entries\n", len(cachedFetcher.cache))
	for body := range control.results {
		fmt.Println("Result: ", body)
	}
}
