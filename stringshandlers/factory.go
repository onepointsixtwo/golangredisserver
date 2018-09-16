package stringshandlers

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/router"
)

const (
	PING   = "PING"
	SET    = "SET"
	GET    = "GET"
	GETSET = "GETSET"
	DEL    = "DEL"
	EXISTS = "EXISTS"
	TIME   = "TIME"
	EXPIRE = "EXPIRE"
	TTL    = "TTL"

	OK = "OK"
)

type Factory struct {
	dataStore     keyvaluestore.Store
	expiryHandler *expiry.Handler
}

func NewFactory(dataStore keyvaluestore.Store, expiryHandler *expiry.Handler) *Factory {
	return &Factory{dataStore, expiryHandler}
}

func (factory *Factory) AddHandlersToRouter(router router.Router) {
	router.AddRedisCommandHandler(PING, factory.handlePing)
	router.AddRedisCommandHandler(GET, factory.handleGet)
	router.AddRedisCommandHandler(SET, factory.handleSet)
	router.AddRedisCommandHandler(GETSET, factory.handleGetSet)
	router.AddRedisCommandHandler(DEL, factory.handleDel)
	router.AddRedisCommandHandler(EXISTS, factory.handleExists)
	router.AddRedisCommandHandler(TIME, factory.handleTime)
	router.AddRedisCommandHandler(EXPIRE, factory.handleExpire)
	router.AddRedisCommandHandler(TTL, factory.handleTTL)
}

func (factory *Factory) handlePing(args []string, connection connection.Connection) {
	pingHandler(args, connection)
}

func (factory *Factory) handleSet(args []string, connection connection.Connection) {
	setHandler(args, connection, factory.dataStore)
}

func (factory *Factory) handleGet(args []string, connection connection.Connection) {
	getHandler(args, connection, factory.dataStore)
}

func (factory *Factory) handleGetSet(args []string, connection connection.Connection) {
	getSetHandler(args, connection, factory.dataStore)
}

func (factory *Factory) handleDel(args []string, connection connection.Connection) {
	deleteHandler(args, connection, factory.dataStore, factory.expiryHandler)
}

func (factory *Factory) handleExists(args []string, connection connection.Connection) {
	existsHandler(args, connection, factory.dataStore)
}

func (factory *Factory) handleTime(args []string, connection connection.Connection) {
	timeHandler(args, connection)
}

func (factory *Factory) handleExpire(args []string, connection connection.Connection) {
	expireHandler(args, connection, factory.expiryHandler)
}

func (factory *Factory) handleTTL(args []string, connection connection.Connection) {
	ttlHandler(args, connection, factory.expiryHandler)
}
