package clientconnection

/*
	Changed my mind about how this class works. It should not be doing anything to do with actually processing commands.
	It's caller should do that, because the outer server class (or whatever) has access to the data structures etc. it
	will need for that. It'll have to be thread safe access of those anyway.
	This file has one purpose - reading data from the connection, and writing data to the connection. It's caller can either call
	Start() for single-threaded use, or go Start() for a thread-per-connection.
*/

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/router"
	"io"
	"net"
	"strings"
)

const (
	CRLF = "\r\n"
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
	conn := connection.connection

	// Create the buffer and the line ending we're looking for for command ends.
	lineEndingBytes := []byte(CRLF)
	buffer := make([]byte, 0)

	for {
		// Read the next 1024 bytes and break on errors.
		readBuffer := make([]byte, 1024)
		bytesRead, err := conn.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				fmt.Errorf("Error reading from connection %v\n", err)
			} else {
				fmt.Printf("Reached EOF from client %v\n", conn.RemoteAddr())
			}
			break
		}

		// Copy bytes across to the buffer, and if we find a line ending, skip the next
		// byte (i.e. \n) and handle the command that's incoming. Also remake the buffer
		// awaiting the next command's ending.
		for i := 0; i < bytesRead; i++ {
			b := readBuffer[i]
			if b == lineEndingBytes[0] {
				command := string(buffer)
				connection.handleCommand(command)
				buffer = make([]byte, 0)
				i++
			} else {
				buffer = append(buffer, b)
			}
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
func (connection *ClientConnection) handleCommand(command string) {
	commandParts := strings.Split(command, " ")

	cmd := commandParts[0]
	args := commandParts[1:]

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
