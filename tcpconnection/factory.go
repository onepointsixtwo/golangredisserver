package tcpconnection

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

type TCPConnectionFactory struct {
	router router.Router
}

func NewConnectionFactory(router router.Router) *TCPConnectionFactory {
	return &TCPConnectionFactory{router: router}
}

func (factory *TCPConnectionFactory) CreateConnection(networkConnection net.Conn, finishedChannel chan<- connection.Connection) connection.Connection {
	return New(networkConnection, factory.router, finishedChannel)
}
