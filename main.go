package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
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

func (s *simpleServer) Address() string {
	return s.address
}

func (s *simpleServer) isAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
} 

func (lb *loadbalancer) getNextAvailableServer() Server {
	// round robin algo
	server := lb.servers[lb.roundRobinCount % len(lb.servers)]
	for !server.isAlive(){
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount % len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *loadbalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, req)
}

func main() {
	servers := []Server{
		newSimpleServer("http://www.twitter.com"),
		newSimpleServer("http://www.amazon.com"),
		newSimpleServer("http://www.google.com"),
	}

	lb := newLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)

	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":" + lb.port, nil)


}