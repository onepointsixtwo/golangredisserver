package mocks

import (
	"github.com/onepointsixtwo/golangredisserver/router"
)

type MockRouter struct {
	CommandsReceived []*ReceivedCommand
}

type ReceivedCommand struct {
	Command   string
	Args      []string
	Responder router.Responder
}

func NewMockRouter() *MockRouter {
	return &MockRouter{make([]*ReceivedCommand, 0)}
}

func (mockRouter *MockRouter) RouteIncomingCommand(command string, args []string, responder router.Responder) error {
	receivedCommand := &ReceivedCommand{command, args, responder}
	mockRouter.CommandsReceived = append(mockRouter.CommandsReceived, receivedCommand)
	return nil
}
