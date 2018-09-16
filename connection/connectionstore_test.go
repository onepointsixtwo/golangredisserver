package connection

import (
	"testing"
)

// Tests

func TestInsertingIntoClientConnectionStore(t *testing.T) {
	store := NewStore()

	connection1 := newMockConnection()
	store.AddClientConnection(connection1)

	if store.GetClientConnectionsCount() != 1 {
		t.Errorf("After inserting 1 connection the store should contain 1 but contains %v", store.GetClientConnectionsCount())
	}

	connection2 := newMockConnection()
	store.AddClientConnection(connection2)

	if store.GetClientConnectionsCount() != 2 {
		t.Errorf("After inserting 2 connections the store should contain 2 but contains %v", store.GetClientConnectionsCount())
	}

	for x := 0; x < 30; x++ {
		store.AddClientConnection(newMockConnection())
	}

	if store.GetClientConnectionsCount() != 32 {
		t.Errorf("After inserting 32 connections in total the store should contain 32 but contains %v", store.GetClientConnectionsCount())
	}
}

func TestDeletingFromClientConnectionStore(t *testing.T) {
	store := NewStore()

	connection1 := newMockConnection()
	store.AddClientConnection(connection1)

	connection2 := newMockConnection()
	store.AddClientConnection(connection2)

	if store.GetClientConnectionsCount() != 2 {
		t.Errorf("After inserting 2 connections the store should contain 2 but contains %v", store.GetClientConnectionsCount())
	}

	store.RemoveClientConnection(connection2)

	if store.GetClientConnectionsCount() != 1 {
		t.Errorf("After deleting an item from a store of 2, the store should contain 1 but contains %v", store.GetClientConnectionsCount())
	}

	store.RemoveClientConnection(connection1)

	if store.GetClientConnectionsCount() != 0 {
		t.Errorf("The store should contain zero connections but contains %v", store.GetClientConnectionsCount())
	}
}

// Mock Connection
type MockConnection struct{}

func newMockConnection() *MockConnection {
	return &MockConnection{}
}

func (conn *MockConnection) Start() {}

func (conn *MockConnection) SendResponse(response string) {}

func (conn *MockConnection) CreateResponseWriter() ConnectionResponseWriter {
	return nil
}
