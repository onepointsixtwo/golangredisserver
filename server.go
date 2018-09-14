package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/responsewriter"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

const (
	PING   = "PING"
	PONG   = "PONG"
	SET    = "SET"
	GET    = "GET"
	DEL    = "DEL"
	EXISTS = "EXISTS"
	OK     = "OK"
	CRLF   = "\r\n"
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
	router.AddRedisCommandHandler(DEL, server.deleteHandler)
	router.AddRedisCommandHandler(EXISTS, server.existsHandler)

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
	writer := responsewriter.New(responder)
	if len(args) > 0 {
		writer.AddBulkString(args[0])
	} else {
		writer.AddSimpleString(PONG)
	}
	server.writeResponse(writer)
}

func (server *RedisServer) getHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 0 {
		key := args[0]
		value, err := server.dataStore.StringForKey(key)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("value not found for key '%v'", key))
		} else {
			writer.AddBulkString(value)
		}
	} else {
		writer.AddErrorString("wrong number of arguments for 'get' command")
	}
	server.writeResponse(writer)
}

func (server *RedisServer) setHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 1 {
		key := args[0]
		value := args[1]
		server.dataStore.SetString(key, value)

		writer.AddSimpleString(OK)
	} else {
		writer.AddErrorString("wrong number of arguments for 'set' command")
	}
	server.writeResponse(writer)
}

func (server *RedisServer) deleteHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)

	deleted := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		success := server.dataStore.DeleteString(key)
		if success {
			deleted++
		}
	}

	writer.AddInt(deleted)
	server.writeResponse(writer)
}

func (server *RedisServer) existsHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)

	exists := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		_, err := server.dataStore.StringForKey(key)
		if err == nil {
			exists++
		}
	}

	writer.AddInt(exists)
	server.writeResponse(writer)
}

// Response helpers

func (server *RedisServer) writeResponse(writer *responsewriter.ResponseWriter) {
	err := writer.WriteResponse()
	if err != nil {
		fmt.Printf("Error writing response: %v", err)
	}
}
