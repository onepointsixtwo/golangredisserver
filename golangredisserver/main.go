package main

/*
	NEXT STEPS

	- Create a router struct similar to gorrilla mux which basically stores a hashmap of functions for each command within Redis and forwards an incoming line to a handler func.
	- Pull out all the logic you can into something separate from the specifics of connection etc.
*/

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver"
	"io"
	"net"
	"os"
	"strings"
)

const (
	TCP  = "tcp"
	CRLF = "\r\n"
)

func handleConnection(connection net.Conn) {
	fmt.Printf("Handling incoming connection from %v\n", connection.RemoteAddr())

	// Create the buffer and the line ending we're looking for for command ends.
	lineEndingBytes := []byte(CRLF)
	buffer := make([]byte, 0)

	for {
		// Read the next 1024 bytes and break on errors.
		readBuffer := make([]byte, 1024)
		bytesRead, err := connection.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				fmt.Errorf("Error reading from connection %v\n", err)
			} else {
				fmt.Print("Reached EOF\n")
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
				handleCommand(connection, command)
				buffer = make([]byte, 0)
				i++
			} else {
				buffer = append(buffer, b)
			}
		}
	}

	// Try to close the connection
	err := connection.Close()
	if err != nil {
		fmt.Printf("Error while attempting to close connection to client %v", err)
	}
}

func handleCommand(connection net.Conn, command string) {
	commandParts := strings.Split(command, " ")
	args := commandParts[1:]

	// So far only handles PING commands. Will add others over time, but after changing the program structure completely.
	switch commandParts[0] {
	case "PING":
		handlePing(connection, args)
	default:
		fmt.Printf("Unknown command %v\n", commandParts[0])
	}
}

func handlePing(connection net.Conn, args []string) {
	fmt.Printf("Handling PING command with args %v\n", args)

	// PING either sends back a pong or the string sent as an argument (if exists)
	if len(args) > 0 {
		fmt.Fprintf(connection, "+%v%v", args[0], CRLF)
	} else {
		fmt.Fprintf(connection, "+PONG%v", CRLF)
	}
}

func main() {
	fmt.Println(golangredisserver.GetString())

	// Golang has lovely networking. This is pretty fun compared to C. Also the docs are actually not shit.

	// Create a listening socket. Kinda weird that the socket is a string, but whatever.
	listeningSocket, err := net.Listen(TCP, ":6379")
	if err != nil {
		fmt.Errorf("Error creating initial listening socket %v\n", err)
		os.Exit(1)
	}

	// Constantly loop awaiting incoming connections
	for {
		connectionToClient, err := listeningSocket.Accept()
		if err != nil {
			fmt.Errorf("Error accepting incoming connection %v\n", err)
		}

		// Incoming connections are handled in a new goro.
		go handleConnection(connectionToClient)
	}
}
