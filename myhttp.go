package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	// BENCHMARK: start chrono.
	start := time.Now()

	parallel := flag.Int("parallel", 10, "max number of parallel requests")
	flag.Parse()

	// Make a channel for the URLs pool.
	urls := make(chan string)

	// The fetch workers are fed with URLs pulled from a channel.
	go func() {
		for _, u := range flag.Args() {
			urls <- u
		}
		close(urls)
	}()

	// The fetch goroutines will write the results in a channel too.
	res := make(chan string)

	// Let the given number of fetch goroutines start working in parallel.
	// We use a waiting group to know when we can close the results channel.
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// To every worker: do fetch urls until the pool is empty.
			for url := range urls {
				fetch(url, res)
			}
		}()
	}

	// Close results channel, when every fetch worker is done.
	go func() {
		wg.Wait()
		close(res)
	}()

	// Pull the results and print.
	// The loop will stop when the result channel is closed.
	for r := range res {
		fmt.Println(r)
	}

	// BENCHMARK: stop and print chrono.
	fmt.Printf("%.2fs total elapsed\n", time.Since(start).Seconds())
}

// fetch writes the hashed response body into a channel, following a request.
func fetch(url string, res chan string) {
	resp, err := http.Get(url)
	if err != nil {
		res <- fmt.Sprint(err)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil { // Write the error to the channel.
		res <- fmt.Sprintf("fetch: reading %s %v", url, err)
		return
	}

	// Write the MD5 hash of the response body to the results channel
	res <- fmt.Sprintf("%s\t%x\n", url, md5.Sum(b))
}
