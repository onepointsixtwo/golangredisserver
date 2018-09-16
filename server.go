package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/router"
	"github.com/onepointsixtwo/golangredisserver/stringshandlers"
	"net"
)

// Server struct

// [ Will also require a key value store to support more than just PING! ]
type RedisServer struct {
	listener                   net.Listener
	router                     router.Router
	connectionCompletedChannel chan *clientconnection.ClientConnection
	connections                *clientconnection.Store
	// This structure will probably change to instead be an array of interface.Factory where interface.Factory is implemented by stringshandler.Factory
	// and any others added later.
	stringHandlerFactory *stringshandlers.Factory
}

// Initialisation

func New(listener net.Listener, dataStore keyvaluestore.Store) *RedisServer {
	connections := clientconnection.NewStore()
	expiryHandler := expiry.New(dataStore)

	stringHandlersFactory := stringshandlers.NewFactory(dataStore, expiryHandler)

	return &RedisServer{listener: listener,
		connections:          connections,
		stringHandlerFactory: stringHandlersFactory}
}

func (server *RedisServer) Init() {
	router := router.NewRedisRouter()
	server.stringHandlerFactory.AddHandlersToRouter(router)
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
