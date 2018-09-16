package connection

import (
	"net"
)

type ConnectionFactory interface {
	CreateConnection(net.Conn, chan<- Connection) Connection
}
