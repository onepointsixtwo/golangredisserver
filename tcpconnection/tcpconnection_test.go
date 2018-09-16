package tcpconnection

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
)

/*
	The same command (PING) is used throughout but it doesn't really matter - this does not handle
	what to do for each incoming connection type. It simply reads from the connection and writes to it.
	Routing handling is delegated to the router interface.
*/

// Tests

func TestTCPConnectionReadsFromConn(t *testing.T) {
	commands := "*2\r\n$4\r\nPING\r\n$4\r\nblah\r\n*1\r\n$4\r\nPING\r\n"
	finishedChannel := make(chan connection.Connection, 1)

	sut, _, mockRouter := createTCPConnectionAndDependencies(commands, finishedChannel)

	sut.Start()

	if len(mockRouter.CommandsReceived) != 2 {
		t.Errorf("Expected mock router to contain 2 received commands but contains %v", len(mockRouter.CommandsReceived))
	}

	first := mockRouter.CommandsReceived[0]
	if first.Command != "PING" || first.Args[0] != "blah" {
		t.Errorf("Expected mock router to have received the first command with type PING but was %v and argument blah but was %v", first.Command, first.Args[0])
	}

	second := mockRouter.CommandsReceived[1]
	if second.Command != "PING" || len(second.Args) > 0 {
		t.Errorf("Expected mock router to have received the second command with type PING but was %v and with no arguments but had %v", second.Command, len(second.Args))
	}
}

func TestTCPConnectionWritesToConn(t *testing.T) {
	command := "*1\r\n$4\r\nPING\r\n"
	finishedChannel := make(chan connection.Connection, 1)

	sut, connection, mockRouter := createTCPConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	response := "+PONG"
	con := mockRouter.CommandsReceived[0].Connection
	con.SendResponse(response)

	if connection.WriteBuffer.String() != response {
		t.Errorf("Error using responder to write back to connection - expected response to be written of %v but was %v", response, connection.WriteBuffer.String())
	}
}

func TestTCPConnectionClosesConnectionWhenReadingComplete(t *testing.T) {
	command := "*1\r\n$4\r\nPING\r\n"
	finishedChannel := make(chan connection.Connection, 1)

	sut, connection, _ := createTCPConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	if connection.Closed == false {
		t.Error("Connection should be closed after processing is complete")
	}
}

func TestTCPConnectionSendsCloseToChannelWhenComplete(t *testing.T) {
	command := "*1\r\n$4\r\nPING\r\n"
	finishedChannel := make(chan connection.Connection)

	sut, _, _ := createTCPConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	closedTCPConnection := <-finishedChannel

	if closedTCPConnection != sut {
		t.Error("When client connection is completed it should send itself as closed to the finished channel")
	}
}

// Helpers - create the client and the mock dependencies needed for all tests
func createTCPConnectionAndDependencies(clientCommands string, finishedChannel chan connection.Connection) (*TCPConnection, *mocks.MockConnection, *mocks.MockRouter) {
	conn, router := createTestDependencies(clientCommands)
	return New(conn, router, finishedChannel), conn, router
}

func createTestDependencies(clientCommands string) (*mocks.MockConnection, *mocks.MockRouter) {
	// Create the mock connection
	mockConnection := mocks.NewMockConnection(clientCommands)

	// Create the mock router
	mockRouter := mocks.NewMockRouter()

	return mockConnection, mockRouter
}
