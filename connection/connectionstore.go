package connection

import (
	"sync"
)

type ConnectionStore struct {
	connectionsMutex *sync.Mutex
	connections      []Connection
}

func NewStore() Store {
	return &ConnectionStore{&sync.Mutex{}, make([]Connection, 0)}
}

func (store *ConnectionStore) AddClientConnection(connection Connection) {
	store.connectionsMutex.Lock()
	defer store.connectionsMutex.Unlock()

	store.connections = append(store.connections, connection)
}

func (store *ConnectionStore) RemoveClientConnection(connection Connection) {
	store.connectionsMutex.Lock()
	defer store.connectionsMutex.Unlock()

	var index = -1
	for i, c := range store.connections {
		if c == connection {
			index = i
			break
		}
	}

	if index >= 0 {
		connectionsLength := len(store.connections)
		store.connections[index] = store.connections[connectionsLength-1]
		store.connections = store.connections[:connectionsLength-1]
	}
}

func (store *ConnectionStore) GetClientConnectionsCount() int {
	store.connectionsMutex.Lock()
	defer store.connectionsMutex.Unlock()
	return len(store.connections)
}
