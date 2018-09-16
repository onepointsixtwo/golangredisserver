package connection

type Store interface {
	AddClientConnection(connection Connection)
	RemoveClientConnection(connection Connection)
	GetClientConnectionsCount() int
}
