package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

const (
	TCP  = "tcp"
	PING = "PING"
	CRLF = "\r\n"
)

// Server struct

type RedisServer struct {
	router                     router.Router
	connectionCompletedChannel chan *clientconnection.ClientConnection
}

// Initialisation

func New() *RedisServer {
	return &RedisServer{}
}

func (server *RedisServer) Init() {
	// Setup the router. It will call back on other goroutines so this whole class should use thread safe access.
	router := router.NewRedisRouter()

	// Only one handler for PING so far... more to come!
	router.AddRedisCommandHandler(PING, server.PingHandler)

	server.router = router

	// Setup the channel for listening if connections to clients are completed
	server.connectionCompletedChannel = make(chan *clientconnection.ClientConnection)

	// TODO: should also initialise some kind of storage for the keys and values when
	// we're supporting more than just PING...
}

func (server *RedisServer) Start() error {
	// Create a listening socket
	//TODO: Make the listening port configurable.
	listeningSocket, err := net.Listen(TCP, ":6379")
	if err != nil {
		return fmt.Errorf("Error creating initial listening socket %v\n", err)
	}

	// Start awaiting completed connections in a goro
	go server.handleCompletedConnections()

	// Constantly loop awaiting incoming connections
	for {
		connectionToClient, err := listeningSocket.Accept()
		if err != nil {
			fmt.Errorf("Error accepting incoming connection %v\n", err)
		}

		clientConn := clientconnection.New(connectionToClient, server.router, server.connectionCompletedChannel)
		server.addClientConnection(clientConn)
		clientConn.Start()
	}
}

// Connections Completed Handling
func (server *RedisServer) handleCompletedConnections() {
	//TODO: listen on connectionCompletedChannel and remove completed connections
}

// Client connection caching
func (server *RedisServer) addClientConnection(connection *clientconnection.ClientConnection) {
	//TODO:
}

func (server *RedisServer) removeClientConnection(connection *clientconnection.ClientConnection) {
	//TODO:
}

// Routing handlers

func (server *RedisServer) PingHandler(args []string, responder router.Responder) {
	// PING either sends back a pong or the string sent as an argument (if exists)
	var response string
	if len(args) > 0 {
		response = fmt.Sprintf("+%v%v", args[0], CRLF)
	} else {
		response = fmt.Sprintf("+PONG%v", CRLF)
	}

	responder.SendResponse(response)
}
