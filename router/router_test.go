package router

import (
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

	_ = router.routeIncomingCommand("PING", []string{})
	_ = router.routeIncomingCommand("PING", []string{})

	_ = router.routeIncomingCommand("GET", []string{})

	if tester.pingCount != 2 {
		t.Errorf("Router did not route PING commands. Expected 2 but got %v", tester.pingCount)
	}

	if tester.getCount != 1 {
		t.Errorf("Router did not route GET command. Expected 1 but got %v", tester.getCount)
	}
}

func TestFailsRoutingToUnknown(t *testing.T) {
	router := NewRedisRouter()

	err := router.routeIncomingCommand("UNKNOWN", []string{})

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

func (routerTester *RedisRouterTester) handlePingCommand(command string, args []string) {
	routerTester.pingCount = routerTester.pingCount + 1
}

func (routerTester *RedisRouterTester) handleGetCommand(command string, args []string) {
	routerTester.getCount = routerTester.getCount + 1
}
