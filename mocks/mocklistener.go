package mocks

import (
	"bytes"
	"io"
	"net"
)

type MockListener struct {
	commands          string
	connectionCreated bool
	Connection        *MockConnection
}

func NewMockListener(commands string) *MockListener {
	return &MockListener{commands: commands}
}

func (listener *MockListener) Accept() (net.Conn, error) {
	if !listener.connectionCreated {
		listener.connectionCreated = true

		mockConnection := &MockConnection{Closed: false, ReadString: listener.commands, WriteBuffer: bytes.NewBufferString("")}
		listener.Connection = mockConnection
		return mockConnection, nil
	}
	return nil, io.EOF
}

func (listener *MockListener) Close() error {
	return nil
}

func (listener *MockListener) Addr() net.Addr {
	return &MockAddr{}
}
