package router

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"strings"
)

// The router interface - allows other layers using router to be testable easily.
type Router interface {
	RouteIncomingCommand(command string, args []string, connection connection.Connection) error
	AddRedisCommandHandler(command string, handler RoutingHandler)
}

// RedisRouter types. The router holds a list of handler functions mapped by string
// which represents the command
type RoutingHandler func([]string, connection.Connection)

type RedisRouter struct {
	handlers map[string]RoutingHandler
}

// RedisRouter functions
func NewRedisRouter() *RedisRouter {
	handlersMap := make(map[string]RoutingHandler)
	return &RedisRouter{handlersMap}
}

func (router *RedisRouter) RouteIncomingCommand(command string, args []string, connection connection.Connection) error {
	handler, found := router.handlers[strings.ToUpper(command)]
	if !found {
		return fmt.Errorf("Unable to find handler for command %v in handlers:\n%v", command, router.handlers)
	}

	handler(args, connection)
	return nil
}

func (router *RedisRouter) AddRedisCommandHandler(command string, handler RoutingHandler) {
	router.handlers[strings.ToUpper(command)] = handler
}
