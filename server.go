package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
	"sync"
)

const (
	TCP  = "tcp"
	PING = "PING"
	CRLF = "\r\n"
)

// Server struct

//TODO: Changed my mind about making listening port configurable. It shouldn't even be created here.
// The actual listener should be created externally and given to the struct as part of its initialiser.
// No point in creating this without it and it makes this class testable!
type RedisServer struct {
	router                     router.Router
	connectionCompletedChannel chan *clientconnection.ClientConnection
	connectionsMutex           *sync.Mutex
	connections                []*clientconnection.ClientConnection
}

// Initialisation

func New() *RedisServer {
	return &RedisServer{connectionsMutex: &sync.Mutex{}, connections: make([]*clientconnection.ClientConnection, 0)}
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
		go clientConn.Start()
	}
}

// Connections Completed Handling
func (server *RedisServer) handleCompletedConnections() {
	for completedClientConnection := range server.connectionCompletedChannel {
		server.removeClientConnection(completedClientConnection)
	}
}

// Client connection caching
//TODO: pull this functionality out into another struct for this purpose.
func (server *RedisServer) addClientConnection(connection *clientconnection.ClientConnection) {
	server.connectionsMutex.Lock()
	defer server.connectionsMutex.Unlock()

	server.connections = append(server.connections, connection)

	fmt.Printf("There are %v client connections\n", len(server.connections))
}

func (server *RedisServer) removeClientConnection(connection *clientconnection.ClientConnection) {
	server.connectionsMutex.Lock()
	defer server.connectionsMutex.Unlock()

	var index = -1
	for i, c := range server.connections {
		if c == connection {
			index = i
			break
		}
	}

	if index >= 0 {
		connectionsLength := len(server.connections)
		server.connections[index] = server.connections[connectionsLength-1]
		server.connections = server.connections[:connectionsLength-1]
	}

	fmt.Printf("There are %v client connections\n", len(server.connections))
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
