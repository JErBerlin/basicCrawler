package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	fmt.Println("Starting tests..")
	var tests = []struct {
		desc string
		url string
		want string
	}{
		{"ok req", "/path1", "8b04d5e3775d298e78455efc5ca404d5"},
		{"ok req", "/path2", "a9f0e61a137d86aa9db53465e0801612"},
	}

	res := make(chan string, len(tests))
	defer close(res)
	cl := http.Client{
		Timeout: time.Second * sTimeOutClient,
	}
	// Test the fetch function to write in the response channel but not in parallel
	for _, test := range tests {
		srv := serverMock(test.url)
		defer srv.Close()

		fmt.Printf("Test run: %v\n", test)
		left := fmt.Sprintf("Fetch(%v)", test.url) // srv.URL
		Fetch(cl, srv.URL+test.url, res)
		got := <-res
		want := srv.URL+test.url+" "+test.want
		if got != want {
			t.Errorf("%s: %s = %q, want %q", test.desc, left, got, want)
		}
	}
}

func serverMock(url string) *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc(url, responseMock)
	srv := httptest.NewServer(handler)

	return srv
}

func responseMock(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/path1":
		_, _ = w.Write([]byte("first"))
	case "/path2":
		_, _ = w.Write([]byte("second"))
	default:
		_, _ = w.Write([]byte("default"))
	}

}
