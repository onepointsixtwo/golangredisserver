package main

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"github.com/onepointsixtwo/golangredisserver/handlers"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/stringshandlers"
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

	// Create other dependencies for handler factories
	keyValueDataStore := keyvaluestore.New()
	expiryHandler := expiry.New(keyValueDataStore)
	connectionsStore := clientconnection.NewStore()

	// Create handler factories
	strings := stringshandlers.NewFactory(keyValueDataStore, expiryHandler)
	factories := []handlers.Factory{strings}

	// Initialise golang redis server with created dependencies
	server := golangredisserver.New(listener, connectionsStore, factories)
	server.Init()

	err = server.Start()
	if err != nil {
		fmt.Errorf("Fatal error - unable to start server %v", err)
	}
}
