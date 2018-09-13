package mocks

import (
	"bytes"
	"net"
	"time"
)

type MockConnection struct {
	Closed        bool
	ReadDeadline  time.Time
	WriteDeadline time.Time
	WriteBuffer   *bytes.Buffer
	reader        *MockReader
}

func NewMockConnection(incomingReadString string) *MockConnection {
	reader := NewMockReader(incomingReadString)
	return &MockConnection{Closed: false, WriteBuffer: &bytes.Buffer{}, reader: reader}
}

func (mockConn *MockConnection) Read(b []byte) (n int, err error) {
	return mockConn.reader.Read(b)
}

func (mockConn *MockConnection) Write(b []byte) (n int, err error) {
	mockConn.WriteBuffer.Write(b)
	return len(b), nil
}

func (mockConn *MockConnection) Close() error {
	mockConn.Closed = true
	return nil
}

func (mockConn *MockConnection) LocalAddr() net.Addr {
	return NewMockAddr("tcp", "0.0.0.0")
}

func (mockConn *MockConnection) RemoteAddr() net.Addr {
	return NewMockAddr("tcp", "0.0.0.0")
}

func (mockConn *MockConnection) SetDeadline(t time.Time) error {
	mockConn.SetReadDeadline(t)
	mockConn.SetWriteDeadline(t)
	return nil
}

func (mockConn *MockConnection) SetReadDeadline(t time.Time) error {
	mockConn.ReadDeadline = t
	return nil
}

func (mockConn *MockConnection) SetWriteDeadline(t time.Time) error {
	mockConn.WriteDeadline = t
	return nil
}
