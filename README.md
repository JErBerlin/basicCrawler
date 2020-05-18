# Adjust Basic Crawler

A tool that makes http requests and prints the address of the request along with the MD5 hash of the response

## How to run..

###  the tests

`go test`

###  the compiler 

`go build myhttp.go`

###  the program 

To fetch just one url

`./myhttp http://www.mysite.com`

fecthing several urls (by default in parallel with 10 threads)

`./myhttp http://www.mysite1.com http://www.mysite2.com`

specifying the number n = 2 of parallel threads (goroutines)

`./myhttp -parallel=2 http://www.mysite1.com http://www.mysite2.com`

taking the urls from a file (text file with one ulr per line), in linux:

`xargs ./myhttp -parallel=100 < urlslist.txt`

## TODO

- Write edge cases tests
