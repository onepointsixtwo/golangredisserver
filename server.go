package golangredisserver

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/clientconnection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
	"github.com/onepointsixtwo/golangredisserver/responsewriter"
	"github.com/onepointsixtwo/golangredisserver/router"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	PING   = "PING"
	PONG   = "PONG"
	SET    = "SET"
	GET    = "GET"
	GETSET = "GETSET"
	DEL    = "DEL"
	EXISTS = "EXISTS"
	TIME   = "TIME"
	EXPIRE = "EXPIRE"
	OK     = "OK"
	CRLF   = "\r\n"
)

// Server struct

// [ Will also require a key value store to support more than just PING! ]
type RedisServer struct {
	listener                   net.Listener
	router                     router.Router
	connectionCompletedChannel chan *clientconnection.ClientConnection
	connections                *clientconnection.Store
	dataStore                  keyvaluestore.Store
	timersMap                  map[string]*time.Timer
	timersMapLock              *sync.Mutex
}

// Initialisation

func New(listener net.Listener) *RedisServer {
	return &RedisServer{listener: listener,
		connections:   clientconnection.NewStore(),
		dataStore:     keyvaluestore.New(),
		timersMap:     make(map[string]*time.Timer),
		timersMapLock: &sync.Mutex{}}
}

func (server *RedisServer) Init() {
	// Setup the router. It will call back on other goroutines so this whole class should use thread safe access.
	router := router.NewRedisRouter()

	// Only one handler for PING so far... more to come!
	router.AddRedisCommandHandler(PING, server.pingHandler)
	router.AddRedisCommandHandler(GET, server.getHandler)
	router.AddRedisCommandHandler(SET, server.setHandler)
	router.AddRedisCommandHandler(GETSET, server.getSetHandler)
	router.AddRedisCommandHandler(DEL, server.deleteHandler)
	router.AddRedisCommandHandler(EXISTS, server.existsHandler)
	router.AddRedisCommandHandler(TIME, server.timeHandler)
	router.AddRedisCommandHandler(EXPIRE, server.expireHandler)

	server.router = router

	// Setup the channel for listening if connections to clients are completed
	server.connectionCompletedChannel = make(chan *clientconnection.ClientConnection)
}

func (server *RedisServer) Start() error {
	defer server.listener.Close()

	go server.handleCompletedConnections()

	for {
		connectionToClient, err := server.listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting incoming connection %v\n", err)
		}
		server.handleNewClient(connectionToClient)
	}
}

// Starting connections
func (server *RedisServer) handleNewClient(conn net.Conn) {
	clientConn := clientconnection.New(conn, server.router, server.connectionCompletedChannel)
	server.connections.AddClientConnection(clientConn)
	go clientConn.Start()
}

// Connections Completed Handling
func (server *RedisServer) handleCompletedConnections() {
	for completedClientConnection := range server.connectionCompletedChannel {
		server.connections.RemoveClientConnection(completedClientConnection)
	}
}

// Routing handlers

func (server *RedisServer) pingHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 0 {
		writer.AddBulkString(args[0])
	} else {
		writer.AddSimpleString(PONG)
	}
	server.writeResponse(writer)
}

func (server *RedisServer) getHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 0 {
		key := args[0]
		value, err := server.dataStore.StringForKey(key)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("value not found for key '%v'", key))
		} else {
			writer.AddBulkString(value)
		}
	} else {
		writer.AddErrorString("wrong number of arguments for 'get' command")
	}
	server.writeResponse(writer)
}

func (server *RedisServer) setHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 1 {
		key := args[0]
		value := args[1]
		server.dataStore.SetString(key, value)

		writer.AddSimpleString(OK)
	} else {
		writer.AddErrorString("wrong number of arguments for 'set' command")
	}
	server.writeResponse(writer)
}

func (server *RedisServer) getSetHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) > 1 {
		key := args[0]
		value := args[1]

		existingValue, _ := server.dataStore.StringForKey(key)
		server.dataStore.SetString(key, value)

		writer.AddBulkString(existingValue)
	} else {
		writer.AddErrorString("wrong number of arguments for 'set' command")
	}
	server.writeResponse(writer)
}

func (server *RedisServer) deleteHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)

	deleted := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		success := server.dataStore.DeleteString(key)
		if success {
			server.cancelTimerForKeyIfExists(key)
			deleted++
		}
	}

	writer.AddInt(deleted)
	server.writeResponse(writer)
}

func (server *RedisServer) existsHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)

	exists := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		_, err := server.dataStore.StringForKey(key)
		if err == nil {
			exists++
		}
	}

	writer.AddInt(exists)
	server.writeResponse(writer)
}

func (server *RedisServer) timeHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)

	currentTime := time.Now()

	//Get the seconds
	seconds := currentTime.Unix()
	writer.AddBulkString(fmt.Sprintf("%v", seconds))

	//Get the microseconds
	nanoseconds := currentTime.UnixNano()
	nanosecondsRemainder := nanoseconds % (seconds * int64(time.Nanosecond))
	milliseconds := nanosecondsRemainder / 1000
	writer.AddBulkString(fmt.Sprintf("%v", milliseconds))

	server.writeResponse(writer)
}

func (server *RedisServer) expireHandler(args []string, responder router.Responder) {
	writer := responsewriter.New(responder)
	if len(args) == 2 {
		key := args[0]
		expirySecondsString := args[1]

		expirySeconds, err := strconv.Atoi(expirySecondsString)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("unable to parse argument for expiry time in seconds: %v", err))
		} else {
			_, err := server.dataStore.StringForKey(key)
			if err != nil {
				writer.AddErrorString("cannot set expiry for non existent key!")
			} else {
				server.expireKey(key, expirySeconds)
				writer.AddSimpleString(OK)
			}
		}
	} else {
		writer.AddErrorString("incorrect number of args - should have two, a key and expiry time")
	}
	server.writeResponse(writer)
}

// Timing Handlers (To be moved to its own separated structure to handle expiry)
func (server *RedisServer) expireKey(key string, afterSeconds int) {
	// Cancel existing timer
	server.cancelTimerForKeyIfExists(key)

	// Start new timer
	timer := time.NewTimer(time.Duration(afterSeconds) * time.Second)
	server.storeTimerForKey(timer, key)
	go server.runTimer(timer, key)
}

func (server *RedisServer) runTimer(timer *time.Timer, key string) {
	<-timer.C

	server.removeTimerForKey(timer, key)

	fmt.Printf("Deleting expiring key: %v\n", key)
	server.dataStore.DeleteString(key)
}

func (server *RedisServer) storeTimerForKey(timer *time.Timer, key string) {
	server.timersMapLock.Lock()
	defer server.timersMapLock.Unlock()

	server.timersMap[key] = timer
}

func (server *RedisServer) removeTimerForKey(timer *time.Timer, key string) {
	server.timersMapLock.Lock()
	defer server.timersMapLock.Unlock()

	_, exists := server.timersMap[key]
	if exists {
		delete(server.timersMap, key)
	}
}

func (server *RedisServer) cancelTimerForKeyIfExists(key string) {
	server.timersMapLock.Lock()
	defer server.timersMapLock.Unlock()

	timer, exists := server.timersMap[key]
	if exists {
		stopped := timer.Stop()
		if !stopped {
			fmt.Println("Failed to stop timer!")
		}
		delete(server.timersMap, key)
	}
}

// Response helpers

func (server *RedisServer) writeResponse(writer *responsewriter.ResponseWriter) {
	err := writer.WriteResponse()
	if err != nil {
		fmt.Printf("Error writing response: %v", err)
	}
}
