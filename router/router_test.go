package router

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"testing"
)

// TESTS

func TestRedisRouterCreation(t *testing.T) {
	redisRouter := NewRedisRouter()
	if redisRouter == nil {
		t.Error("NewRedisRouter() should create a new instance of *RedisRouter")
	}
}

func TestSuccessfulRouting(t *testing.T) {
	tester := newRouterTester()

	router := NewRedisRouter()
	router.AddRedisCommandHandler("PING", tester.handlePingCommand)
	router.AddRedisCommandHandler("GET", tester.handleGetCommand)

	_ = router.RouteIncomingCommand("PING", []string{}, nil)
	_ = router.RouteIncomingCommand("PING", []string{}, nil)

	_ = router.RouteIncomingCommand("GET", []string{}, nil)

	if tester.pingCount != 2 {
		t.Errorf("Router did not route PING commands. Expected 2 but got %v", tester.pingCount)
	}

	if tester.getCount != 1 {
		t.Errorf("Router did not route GET command. Expected 1 but got %v", tester.getCount)
	}
}

func TestFailsRoutingToUnknown(t *testing.T) {
	router := NewRedisRouter()

	err := router.RouteIncomingCommand("UNKNOWN", []string{}, nil)

	if err == nil {
		t.Error("Expected error when routing unknown command. Router failed.")
	}
}

// HELPERS

type RedisRouterTester struct {
	pingCount, getCount int
}

func newRouterTester() *RedisRouterTester {
	return &RedisRouterTester{0, 0}
}

func (routerTester *RedisRouterTester) handlePingCommand(args []string, connection connection.Connection) {
	routerTester.pingCount = routerTester.pingCount + 1
}

func (routerTester *RedisRouterTester) handleGetCommand(args []string, connection connection.Connection) {
	routerTester.getCount = routerTester.getCount + 1
}
