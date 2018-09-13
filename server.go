package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

const (
	PING = "PING"
	SET  = "SET"
	GET  = "GET"
	OK   = "OK"
	CRLF = "\r\n"
)

// Server struct

// [ Will also require a key value store to support more than just PING! ]
type RedisServer struct {
	listener                   net.Listener
	router                     router.Router
	connectionCompletedChannel chan *clientconnection.ClientConnection
	connections                *clientconnection.Store
	dataStore                  keyvaluestore.Store
}

// Initialisation

func New(listener net.Listener) *RedisServer {
	return &RedisServer{listener: listener, connections: clientconnection.NewStore(), dataStore: keyvaluestore.New()}
}

func (server *RedisServer) Init() {
	// Setup the router. It will call back on other goroutines so this whole class should use thread safe access.
	router := router.NewRedisRouter()

	// Only one handler for PING so far... more to come!
	router.AddRedisCommandHandler(PING, server.pingHandler)
	router.AddRedisCommandHandler(GET, server.getHandler)
	router.AddRedisCommandHandler(SET, server.setHandler)

	server.router = router

	// Setup the channel for listening if connections to clients are completed
	server.connectionCompletedChannel = make(chan *clientconnection.ClientConnection)
}

func (server *RedisServer) Start() error {
	defer server.listener.Close()

	go server.handleCompletedConnections()

	for {
		connectionToClient, err := server.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting incoming connection %v\n", err)
		}
		server.handleNewClient(connectionToClient)
	}
}

// Starting connections
func (server *RedisServer) handleNewClient(conn net.Conn) {
	clientConn := clientconnection.New(conn, server.router, server.connectionCompletedChannel)
	server.connections.AddClientConnection(clientConn)
	go clientConn.Start()
}

// Connections Completed Handling
func (server *RedisServer) handleCompletedConnections() {
	for completedClientConnection := range server.connectionCompletedChannel {
		server.connections.RemoveClientConnection(completedClientConnection)
	}
}

// Routing handlers

func (server *RedisServer) pingHandler(args []string, responder router.Responder) {
	// PING either sends back a pong or the string sent as an argument (if exists)
	var response string
	if len(args) > 0 {
		response = server.redisBulkStringifyValue(args[0])
	} else {
		response = fmt.Sprintf("+PONG%v", CRLF)
	}

	responder.SendResponse(response)
}

func (server *RedisServer) getHandler(args []string, responder router.Responder) {
	key := args[0]

	if key != "" {
		value, err := server.dataStore.StringForKey(key)
		if err != nil {
			responder.SendResponse(server.errorStringifyValue(fmt.Sprintf("value not found for key '%v'", key)))
		} else {
			responder.SendResponse(server.redisBulkStringifyValue(value))
		}
	} else {
		responder.SendResponse(server.errorStringifyValue("wrong number of arguments for 'get' command"))
	}
}

func (server *RedisServer) setHandler(args []string, responder router.Responder) {
	key := args[0]
	value := args[1]

	if key != "" && value != "" {
		server.dataStore.SetString(key, value)
		responder.SendResponse(fmt.Sprintf("+%v%v", OK, CRLF))
	} else {
		responder.SendResponse(server.errorStringifyValue("wrong number of arguments for 'set' command"))
	}
}

// Response helpers

func (server *RedisServer) redisBulkStringifyValue(value string) string {
	return fmt.Sprintf("$%v%v%v%v", len(value), CRLF, value, CRLF)
}

func (server *RedisServer) errorStringifyValue(errorString string) string {
	return fmt.Sprintf("-%v%v", errorString, CRLF)
}
