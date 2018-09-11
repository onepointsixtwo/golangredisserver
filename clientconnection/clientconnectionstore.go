package clientconnection

import (
	"sync"
)

type Store struct {
	connectionsMutex *sync.Mutex
	connections      []*ClientConnection
}

func NewStore() *Store {
	return &Store{&sync.Mutex{}, make([]*ClientConnection, 0)}
}

func (store *Store) AddClientConnection(connection *ClientConnection) {
	store.connectionsMutex.Lock()
	defer store.connectionsMutex.Unlock()

	store.connections = append(store.connections, connection)
}

func (store *Store) RemoveClientConnection(connection *ClientConnection) {
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

func (store *Store) GetClientConnectionsCount() int {
	store.connectionsMutex.Lock()
	defer store.connectionsMutex.Unlock()
	return len(store.connections)
}
