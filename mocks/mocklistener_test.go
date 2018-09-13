package mocks

import (
	"net"
	"testing"
)

// Tests

func TestMockListenerIsNetListener(t *testing.T) {
	var listener net.Listener
	listener = NewMockListener("PING\r\n")

	if listener == nil {
		t.Error("Listener should be castable to net.Listener")
	}
}

func TestMockListenerAcceptsSingleConnectionWithGivenInputString(t *testing.T) {
	inputString := "PING\r\n"
	listener := NewMockListener(inputString)

	connection, err := listener.Accept()

	if err != nil || connection == nil || readAllConnectionInputToString(connection) != inputString {
		t.Error("Mock listener should accept an initial connection with given input string from its read function")
	}

	connection2, err2 := listener.Accept()

	if err2 == nil || connection2 != nil {
		t.Error("Mock listener should only have one connection to accept and then should error")
	}
}
