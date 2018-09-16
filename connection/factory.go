package connection

import (
	"net"
)

type ConnectionFactory interface {
	CreateConnection(networkConnection net.Conn, finishedChannel chan<- Connection) Connection
}
