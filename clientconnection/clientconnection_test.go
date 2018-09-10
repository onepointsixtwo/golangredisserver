package clientconnection

import (
	"bytes"
	"github.com/onepointsixtwo/golangredisserver/router"
	"io"
	"net"
	"testing"
	"time"
)

/*
	The same command (PING) is used throughout but it doesn't really matter - this does not handle
	what to do for each incoming connection type. It simply reads from the connection and writes to it.
	Routing handling is delegated to the router interface.
*/

// Tests

func TestClientConnectionReadsFromConn(t *testing.T) {
	commands := "PING blah\r\nPING\r\n"
	finishedChannel := make(chan *ClientConnection, 1)

	sut, _, mockRouter := createClientConnectionAndDependencies(commands, finishedChannel)

	sut.Start()

	if len(mockRouter.commandsReceived) != 2 {
		t.Errorf("Expected mock router to contain 2 received commands but contains %v", len(mockRouter.commandsReceived))
	}

	first := mockRouter.commandsReceived[0]
	if first.command != "PING" || first.args[0] != "blah" {
		t.Errorf("Expected mock router to have received the first command with type PING but was %v and argument blah but was %v", first.command, first.args[0])
	}

	second := mockRouter.commandsReceived[1]
	if second.command != "PING" || len(second.args) > 0 {
		t.Errorf("Expected mock router to have received the second command with type PING but was %v and with no arguments but had %v", second.command, len(second.args))
	}
}

func TestClientConnectionWritesToConn(t *testing.T) {
	command := "PING\r\n"
	finishedChannel := make(chan *ClientConnection, 1)

	sut, connection, mockRouter := createClientConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	response := "+PONG"
	responder := mockRouter.commandsReceived[0].responder
	responder.SendResponse(response)

	if connection.writeBuffer.String() != response {
		t.Errorf("Error using responder to write back to connection - expected response to be written of %v but was %v", response, connection.writeBuffer.String())
	}
}

func TestClientConnectionClosesConnectionWhenReadingComplete(t *testing.T) {
	command := "PING\r\n"
	finishedChannel := make(chan *ClientConnection, 1)

	sut, connection, _ := createClientConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	if connection.closed == false {
		t.Error("Connection should be closed after processing is complete")
	}
}

func TestClientConnectionSendsCloseToChannelWhenComplete(t *testing.T) {
	command := "PING\r\n"
	finishedChannel := make(chan *ClientConnection)

	sut, _, _ := createClientConnectionAndDependencies(command, finishedChannel)

	sut.Start()

	closedClientConnection := <-finishedChannel

	if closedClientConnection != sut {
		t.Error("When client connection is completed it should send itself as closed to the finished channel")
	}
}

// Helpers - create the client and the mock dependencies needed for all tests
func createClientConnectionAndDependencies(clientCommands string, finishedChannel chan *ClientConnection) (*ClientConnection, *MockConnection, *MockRouter) {
	conn, router := createTestDependencies(clientCommands)
	return New(conn, router, finishedChannel), conn, router
}

func createTestDependencies(clientCommands string) (*MockConnection, *MockRouter) {
	// Create the mock connection
	mockConnection := &MockConnection{closed: false, readString: clientCommands, writeBuffer: bytes.NewBufferString("")}

	// Create the mock router
	mockRouter := &MockRouter{make([]*ReceivedCommand, 0)}

	return mockConnection, mockRouter
}

// Mocking

type MockConnection struct {
	closed            bool
	readDeadline      time.Time
	writeDeadline     time.Time
	readString        string
	currentReadOffset int
	writeBuffer       *bytes.Buffer
}

func (mockConn *MockConnection) Read(b []byte) (n int, err error) {
	lengthOfGivenArray := len(b)

	// Not hugely efficient to convert every time...
	arrayFromReadString := []byte(mockConn.readString)
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
	mockConn.writeBuffer.Write(b)
	return len(b), nil
}

func (mockConn *MockConnection) Close() error {
	mockConn.closed = true
	return nil
}

func (mockConn *MockConnection) LocalAddr() net.Addr {
	return &MockAddr{}
}

func (mockConn *MockConnection) RemoteAddr() net.Addr {
	return &MockAddr{}
}

func (mockConn *MockConnection) SetDeadline(t time.Time) error {
	mockConn.SetReadDeadline(t)
	mockConn.SetWriteDeadline(t)
	return nil
}

func (mockConn *MockConnection) SetReadDeadline(t time.Time) error {
	mockConn.readDeadline = t
	return nil
}

func (mockConn *MockConnection) SetWriteDeadline(t time.Time) error {
	mockConn.writeDeadline = t
	return nil
}

type MockAddr struct{}

func (addr *MockAddr) Network() string {
	return "tcp"
}

func (addr *MockAddr) String() string {
	return "0.0.0.0"
}

type MockRouter struct {
	commandsReceived []*ReceivedCommand
}

type ReceivedCommand struct {
	command   string
	args      []string
	responder router.Responder
}

func (mockRouter *MockRouter) RouteIncomingCommand(command string, args []string, responder router.Responder) error {
	receivedCommand := &ReceivedCommand{command, args, responder}
	mockRouter.commandsReceived = append(mockRouter.commandsReceived, receivedCommand)
	return nil
}
