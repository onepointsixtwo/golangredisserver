package golangredisserver

import (
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
	"time"
)

// Tests

func TestPingWithoutExtraData(t *testing.T) {
	runServerTest("PING\r\n", nil, func(response string, sut *RedisServer) {
		if response != "+PONG\r\n" {
			t.Errorf("Response to PING should be +PONG\r\n but was %v", response)
		}
	})
}

func TestPingWithExtraData(t *testing.T) {
	runServerTest("PING extra-data\r\n", nil, func(response string, sut *RedisServer) {
		if response != "$10\r\nextra-data\r\n" {
			t.Errorf("Response to PING with arg 'extra-data' should be +extra-data\r\n but was %v", response)
		}
	})
}

func TestSetValueWithGoodKeyAndValue(t *testing.T) {
	runServerTest("SET mykey myvalue\r\n", nil, func(response string, sut *RedisServer) {
		value, _ := sut.dataStore.StringForKey("mykey")
		if value != "myvalue" || response != "+OK\r\n" {
			t.Errorf("Response to SET mykey myvalue should be +OK and value should be in store, but response is %v, value in store is %v", response, value)
		}
	})
}

func TestGetValueWithExistingKey(t *testing.T) {
	store := keyvaluestore.New()
	store.SetString("mykey", "myvalue")

	runServerTest("GET mykey\r\n", store, func(response string, sut *RedisServer) {
		expected := "$7\r\nmyvalue\r\n"
		if response != expected {
			t.Errorf("Response to GET mykey was expected to be %v but was %v", expected, response)
		}
	})
}

// Test Runner

type ServerResponse func(string, *RedisServer)

func runServerTest(clientCommands string, store keyvaluestore.Store, response ServerResponse) {
	// Create the listener
	listener := mocks.NewMockListener(clientCommands)

	// Create the sut (RedisServer) with the created listener.
	sut := New(listener)

	// Replace the store to pre-filled if exists
	if store != nil {
		sut.dataStore = store
	}

	sut.Init()

	go sut.Start()

	// Wait until the server closes the connection - then we're done with processing the single accept connection
	// that mock listener will give back.
	for listener.IsClosed() == false {
		time.Sleep(200 * time.Millisecond)
	}

	output := listener.Connection.WriteBuffer.String()
	response(output, sut)
}
