package mocks

import (
	"bytes"
	"io"
	"net"
	"time"
)

type MockConnection struct {
	Closed            bool
	ReadDeadline      time.Time
	WriteDeadline     time.Time
	ReadString        string
	currentReadOffset int
	WriteBuffer       *bytes.Buffer
}

func NewMockConnection(incomingReadString string) *MockConnection {
	return &MockConnection{Closed: false, ReadString: incomingReadString, WriteBuffer: &bytes.Buffer{}}
}

func (mockConn *MockConnection) Read(b []byte) (n int, err error) {
	lengthOfGivenArray := len(b)

	// Not hugely efficient to convert every time...
	arrayFromReadString := []byte(mockConn.ReadString)
	totalReadStringLength := len(arrayFromReadString)
	remainingReadString := totalReadStringLength - mockConn.currentReadOffset

	if remainingReadString <= 0 {
		return 0, io.EOF
	} else {
		lengthToRead := remainingReadString
		if lengthOfGivenArray < lengthToRead {
			lengthToRead = lengthOfGivenArray
		}

		position := 0
		for i := mockConn.currentReadOffset; i < (mockConn.currentReadOffset + lengthToRead); i++ {
			b[position] = arrayFromReadString[i]
			position++
		}

		mockConn.currentReadOffset += lengthToRead

		return lengthToRead, nil
	}
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
