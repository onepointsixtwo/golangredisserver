package golangredisserver

import (
	"fmt"
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
	connectionFactory          connection.ConnectionFactory
	connectionCompletedChannel chan connection.Connection
}

// Initialisation

func New(listener net.Listener, router router.Router, connections connection.Store, handlerFactories []handlers.Factory, connectionFactory connection.ConnectionFactory) *RedisServer {
	return &RedisServer{listener: listener,
		router:                     router,
		connections:                connections,
		handlerFactories:           handlerFactories,
		connectionFactory:          connectionFactory,
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
		conn, err := server.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting incoming connection %v\n", err)
		}
		server.handleNewClient(conn)
	}
}

// Starting connections
func (server *RedisServer) handleNewClient(conn net.Conn) {
	connection := server.connectionFactory.CreateConnection(conn, server.connectionCompletedChannel)
	server.connections.AddClientConnection(connection)
	go connection.Start()
}

// Connections Completed Handling
func (server *RedisServer) handleCompletedConnections() {
	for completedClientConnection := range server.connectionCompletedChannel {
		server.connections.RemoveClientConnection(completedClientConnection)
	}
}
