package router

import (
	"fmt"
)

// The router interface - allows other layers using router to be testable easily.
type Router interface {
	routeIncomingCommand(command string, args []string) error
}

// RedisRouter types. The router holds a list of handler functions mapped by string
// which represents the command
type RoutingHandler func(string, []string)

type RedisRouter struct {
	handlers map[string]RoutingHandler
}

// RedisRouter functions
func NewRedisRouter() *RedisRouter {
	handlersMap := make(map[string]RoutingHandler)
	return &RedisRouter{handlersMap}
}

func (router *RedisRouter) routeIncomingCommand(command string, args []string) error {
	handler, found := router.handlers[command]
	if !found {
		return fmt.Errorf("Unable to find handler for command %v in handlers:\n%v", command, router.handlers)
	}

	handler(command, args)
	return nil
}

func (router *RedisRouter) AddRedisCommandHandler(command string, handler RoutingHandler) {
	router.handlers[command] = handler
}
