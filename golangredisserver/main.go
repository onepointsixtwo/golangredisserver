package main

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"net"
)

const (
	TCP = "tcp"
)

func main() {
	// Create a listener to instantiate our golang redis server with.
	listener, err := net.Listen(TCP, ":6379")
	if err != nil {
		fmt.Errorf("Error creating initial listening socket %v\n", err)
		return
	}

	server := golangredisserver.New(listener, keyvaluestore.New())
	server.Init()

	err = server.Start()
	if err != nil {
		fmt.Errorf("Fatal error - unable to start server %v", err)
	}
}
