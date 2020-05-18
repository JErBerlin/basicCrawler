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

// Set a reasonable time out (s) for the requests.
const sTimeOutClient = 8

func main() {
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

	// Set a client with time out for all workers.
	// The Client's transport has cached TCP connections, so we should reuse it. It is also safe for concurrent
	// use by multiple goroutines.
	cl := http.Client{
		Timeout: time.Second * sTimeOutClient,
	}

	// Let the given number of fetch goroutines start working in parallel.
	// Use a waiting group to know when we can close the results channel.
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// To every worker: do fetch urls until the pool is empty.
			for url := range urls {
				Fetch(cl, url, res)
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
}

// Fetch writes the hashed response body into a channel, following a request.
func Fetch(cl http.Client, url string, res chan string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		res <- fmt.Sprint(err, "\n")
		return
	}
	resp, err := cl.Do(req)
	if err != nil {
		res <- fmt.Sprint(err, "\n")
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		res <- fmt.Sprint("read:", url, err, "\n")
		return
	}

	// Write the MD5 hash of the response body to the results channel
	res <- fmt.Sprintf("%s %x", url, md5.Sum(b))
}
