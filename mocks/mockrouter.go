package mocks

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
)

type MockRouter struct {
	CommandsReceived []*ReceivedCommand
}

type ReceivedCommand struct {
	Command    string
	Args       []string
	Connection connection.Connection
}

func NewMockRouter() *MockRouter {
	return &MockRouter{make([]*ReceivedCommand, 0)}
}

func (mockRouter *MockRouter) RouteIncomingCommand(command string, args []string, connection connection.Connection) error {
	receivedCommand := &ReceivedCommand{command, args, connection}
	mockRouter.CommandsReceived = append(mockRouter.CommandsReceived, receivedCommand)
	return nil
}
