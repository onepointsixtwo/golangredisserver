package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/handlers"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

// Server struct

// [ Will also require a key value store to support more than just PING! ]
type RedisServer struct {
	listener                   net.Listener
	router                     router.Router
	connections                connection.Store
	handlerFactories           []handlers.Factory
	connectionCompletedChannel chan connection.Connection
}

// Initialisation

func New(listener net.Listener, connections connection.Store, handlerFactories []handlers.Factory) *RedisServer {
	return &RedisServer{listener: listener,
		router:                     router.NewRedisRouter(),
		connections:                connections,
		handlerFactories:           handlerFactories,
		connectionCompletedChannel: make(chan connection.Connection)}
}

func (server *RedisServer) Init() {
	for i := 0; i < len(server.handlerFactories); i++ {
		server.handlerFactories[i].AddHandlersToRouter(server.router)
	}
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
