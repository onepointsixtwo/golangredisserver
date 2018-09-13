package mocks

import (
	"github.com/onepointsixtwo/golangredisserver/router"
	"testing"
)

func TestMockRouterCanCastToRouter(t *testing.T) {
	var router router.Router
	router = NewMockRouter()
	if router == nil {
		t.Error("Mock router should be castable to router.Router")
	}
}

func TestMockRouterRouting(t *testing.T) {
	router := NewMockRouter()

	args := make([]string, 1)
	args[0] = "blah"
	router.RouteIncomingCommand("PING", args, nil)

	if router.CommandsReceived[0] == nil {
		t.Error("Router should have received an incoming command but has none")
	}

	commandReceived := router.CommandsReceived[0]
	if commandReceived.Command != "PING" || commandReceived.Args[0] != "blah" {
		t.Errorf("Incorrect data input to mock router's received commands. Command should have been PING but was %v, args[0] should have been blah but was %v", commandReceived.Command, commandReceived.Args[0])
	}
}
