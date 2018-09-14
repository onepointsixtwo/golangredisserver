package golangredisserver

import (
	"bytes"
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
	"time"
)

// Tests

func TestPingWithoutExtraData(t *testing.T) {
	command := createCommandString("PING")
	runServerTest(command, nil, func(response string, sut *RedisServer) {
		if response != "+PONG\r\n" {
			t.Errorf("Response to PING should be +PONG\r\n but was %v", response)
		}
	})
}

func TestPingWithExtraData(t *testing.T) {
	command := createCommandString("PING", "extra-data")
	runServerTest(command, nil, func(response string, sut *RedisServer) {
		if response != "$10\r\nextra-data\r\n" {
			t.Errorf("Response to PING with arg 'extra-data' should be +extra-data\r\n but was %v", response)
		}
	})
}

func TestSetValueWithGoodKeyAndValue(t *testing.T) {
	command := createCommandString("SET", "mykey", "myvalue")
	runServerTest(command, nil, func(response string, sut *RedisServer) {
		value, _ := sut.dataStore.StringForKey("mykey")
		if value != "myvalue" || response != "+OK\r\n" {
			t.Errorf("Response to SET mykey myvalue should be +OK and value should be in store, but response is %v, value in store is %v", response, value)
		}
	})
}

func TestGetValueWithExistingKey(t *testing.T) {
	store := keyvaluestore.New()
	store.SetString("mykey", "myvalue")

	command := createCommandString("GET", "mykey")

	runServerTest(command, store, func(response string, sut *RedisServer) {
		expected := "$7\r\nmyvalue\r\n"
		if response != expected {
			t.Errorf("Response to GET mykey was expected to be %v but was %v", expected, response)
		}
	})
}

func TestDeleteValueForExistingKey(t *testing.T) {
	store := keyvaluestore.New()
	store.SetString("mykey", "myvalue")
	store.SetString("mykey2", "myvalue2")

	// Delete two keys which exist and attempt one which doesn't. Should give back '2'
	// for those it successfully deleted.
	command := createCommandString("DEL", "mykey", "mykey2", "mykey3")

	runServerTest(command, store, func(response string, sut *RedisServer) {
		expected := ":2\r\n"
		if response != expected {
			t.Errorf("Response to GET mykey was expected to be %v but was %v", expected, response)
		}
	})
}

// Command builder
func createCommandString(command string, args ...string) string {
	var buffer bytes.Buffer

	// Add the length 'header'
	length := 1 + len(args)
	fmt.Fprintf(&buffer, "*%v%v", length, CRLF)

	allStrings := append([]string{command}, args...)
	loops := len(allStrings)
	for i := 0; i < loops; i++ {
		nextStr := allStrings[i]
		strLen := len(nextStr)
		fmt.Fprintf(&buffer, "$%v%v%v%v", strLen, CRLF, nextStr, CRLF)
	}

	return buffer.String()
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
