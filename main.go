package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
)

type simpleServer struct {
	address string
	proxy *httputil.ReverseProxy
}

func newSimpleServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	handleError(err)

	return &simpleServer{
		address: address,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}