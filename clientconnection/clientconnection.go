package clientconnection

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/reader"
	"github.com/onepointsixtwo/golangredisserver/responsewriter"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
	"time"
)

type ClientConnection struct {
	connection      net.Conn
	router          router.Router
	finishedChannel chan<- *ClientConnection
	timeout         time.Duration
}

func New(connection net.Conn, router router.Router, finished chan<- *ClientConnection) *ClientConnection {
	timeout := time.Duration(5) * time.Second
	return &ClientConnection{connection: connection, router: router, finishedChannel: finished, timeout: timeout}
}

func (connection *ClientConnection) Start() {
	fmt.Printf("Handling incoming connection from %v\n", connection.connection.RemoteAddr())

	connection.readAllFromConnection()
	connection.closeConnection()

	go connection.sendConnectionCloseToChannel()
}

// Sends that the connection has been closed to the channel
func (connection *ClientConnection) sendConnectionCloseToChannel() {
	connection.finishedChannel <- connection
}

// Keeps reading from the connection while there's still data to read, and handling
// the incoming commands
func (connection *ClientConnection) readAllFromConnection() {
	readCommand := reader.CreateRespCommandReader(connection.connection)

	for {
		connection.connection.SetReadDeadline(time.Now().Add(connection.timeout))

		command, err := readCommand()

		if err != nil {
			fmt.Errorf("Error trying to read next command %v", err)
			return
		}

		err = connection.handleCommand(command)
		if err != nil {
			writer := responsewriter.New(connection)
			writer.AddErrorString(fmt.Sprintf("unknown command: %v", command.Command))
			_ = writer.WriteResponse()
		}
	}
}

// Closes the connection to the client
func (connection *ClientConnection) closeConnection() {
	err := connection.connection.Close()
	if err != nil {
		fmt.Printf("Error while attempting to close connection to client %v\n", err)
	}
}

// Handles commands read from the incoming connection
func (connection *ClientConnection) handleCommand(respCommand *reader.RespCommand) error {
	cmd := respCommand.Command
	args := respCommand.Args

	fmt.Printf("Handling command from client %v: %v (args: %v)\n", connection.connection.RemoteAddr(), cmd, args)

	return connection.router.RouteIncomingCommand(cmd, args, connection)
}

// connection.Connection implementation

func (connection *ClientConnection) SendResponse(response string) {
	connection.connection.SetWriteDeadline(time.Now().Add(connection.timeout))
	fmt.Fprintf(connection.connection, "%v", response)
}

func (connection *ClientConnection) CreateResponseWriter() connection.ConnectionResponseWriter {
	return responsewriter.New(connection)
}
