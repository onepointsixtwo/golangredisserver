package main

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"github.com/onepointsixtwo/golangredisserver/handlers"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/router"
	"github.com/onepointsixtwo/golangredisserver/stringshandlers"
	"github.com/onepointsixtwo/golangredisserver/tcpconnection"
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
	connectionsStore := connection.NewStore()

	// Create handler factories
	strings := stringshandlers.NewFactory(keyValueDataStore, expiryHandler)
	commandHandlerFactories := []handlers.Factory{strings}

	// Create router
	router := router.NewRedisRouter()

	// Create TCP connection factory
	connectionFactory := tcpconnection.NewConnectionFactory(router)

	// Initialise golang redis server with created dependencies
	server := golangredisserver.New(listener, router, connectionsStore, commandHandlerFactories, connectionFactory)
	server.Init()

	err = server.Start()
	if err != nil {
		fmt.Errorf("Fatal error - unable to start server %v", err)
	}
}
