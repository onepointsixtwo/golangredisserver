package clientconnection

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/reader"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
)

type ClientConnection struct {
	connection      net.Conn
	router          router.Router
	finishedChannel chan<- *ClientConnection
}

func New(connection net.Conn, router router.Router, finished chan<- *ClientConnection) *ClientConnection {
	return &ClientConnection{connection: connection, router: router, finishedChannel: finished}
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
		command, err := readCommand()

		if err != nil {
			fmt.Errorf("Error trying to read next command %v", err)
			return
		}

		connection.handleCommand(command)
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
func (connection *ClientConnection) handleCommand(respCommand *reader.RespCommand) {
	cmd := respCommand.Command
	args := respCommand.Args

	fmt.Printf("Handling command from client %v: %v (args: %v)\n", connection.connection.RemoteAddr(), cmd, args)

	err := connection.router.RouteIncomingCommand(cmd, args, connection)
	if err != nil {
		fmt.Printf("Error routing incoming command %v", err)
	}
}

// Responder implementation - sends response to client from routed command
func (connection *ClientConnection) SendResponse(response string) {
	fmt.Fprintf(connection.connection, "%v", response)
}
