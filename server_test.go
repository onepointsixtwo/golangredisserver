package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
	"time"
)

// Tests

func TestPingWithoutExtraData(t *testing.T) {
	runServerTest("PING\r\n", func(response string) {
		if response != "+PONG\r\n" {
			t.Errorf("Response to PING should be +PONG\r\n but was %v", response)
		}
	})
}

func TestPingWithExtraData(t *testing.T) {
	runServerTest("PING extra-data\r\n", func(response string) {
		if response != "+extra-data\r\n" {
			t.Errorf("Response to PING with arg 'extra-data' should be +extra-data\r\n but was %v", response)
		}
	})
}

// Test Runner

type ServerResponse func(string)

func runServerTest(clientCommands string, response ServerResponse) {
	// Create the listener
	listener := mocks.NewMockListener(clientCommands)

	// Create the sut (RedisServer) with the created listener.
	sut := New(listener)

	sut.Init()

	go sut.Start()

	// Wait until the server closes the connection - then we're done with processing the single accept connection
	// that mock listener will give back.
	for listener.IsClosed() == false {
		time.Sleep(200 * time.Millisecond)
	}

	output := listener.Connection.WriteBuffer.String()
	response(output)
}
