package router

import (
	"fmt"
	"strings"
)

// The router interface - allows other layers using router to be testable easily.
type Responder interface {
	SendResponse(response string)
}

type Router interface {
	RouteIncomingCommand(command string, args []string, responder Responder) error
}

// RedisRouter types. The router holds a list of handler functions mapped by string
// which represents the command
type RoutingHandler func([]string, Responder)

type RedisRouter struct {
	handlers map[string]RoutingHandler
}

// RedisRouter functions
func NewRedisRouter() *RedisRouter {
	handlersMap := make(map[string]RoutingHandler)
	return &RedisRouter{handlersMap}
}

func (router *RedisRouter) RouteIncomingCommand(command string, args []string, responder Responder) error {
	handler, found := router.handlers[strings.ToUpper(command)]
	if !found {
		return fmt.Errorf("Unable to find handler for command %v in handlers:\n%v", command, router.handlers)
	}

	handler(args, responder)
	return nil
}

func (router *RedisRouter) AddRedisCommandHandler(command string, handler RoutingHandler) {
	router.handlers[strings.ToUpper(command)] = handler
}
