package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	address() string
	isAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request) 
}

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

type loadbalancer struct {
	port string
	roundRobinCount int
	servers []Server
}

func newLoadBalancer(port string, servers []Server) *loadbalancer {
	return &loadbalancer{
		port: port,
		servers: servers,
		roundRobinCount: 0,
	}

}

func (lb *loadbalancer) getNextAvailableServer() Server {
	
}

func (lb *loadbalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {

}