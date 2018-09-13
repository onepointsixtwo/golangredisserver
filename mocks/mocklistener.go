package mocks

import (
	"io"
	"net"
	"sync"
)

type MockListener struct {
	commands          string
	connectionCreated bool
	Connection        *MockConnection
	Closed            bool
	closedLock        *sync.Mutex
}

func NewMockListener(commands string) *MockListener {
	return &MockListener{commands: commands, closedLock: &sync.Mutex{}}
}

func (listener *MockListener) Accept() (net.Conn, error) {
	if !listener.connectionCreated {
		listener.connectionCreated = true

		mockConnection := NewMockConnection(listener.commands)
		listener.Connection = mockConnection
		return mockConnection, nil
	}
	return nil, io.EOF
}

func (listener *MockListener) Close() error {
	listener.closedLock.Lock()
	defer listener.closedLock.Unlock()
	listener.Closed = true
	return nil
}

func (listener *MockListener) IsClosed() bool {
	listener.closedLock.Lock()
	defer listener.closedLock.Unlock()
	return listener.Closed
}

func (listener *MockListener) Addr() net.Addr {
	return &MockAddr{}
}
