package stringshandlers

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
)

func getSetHandler(args []string, connection connection.Connection, dataStore keyvaluestore.Store) {
	writer := connection.CreateResponseWriter()
	if len(args) > 1 {
		key := args[0]
		value := args[1]

		existingValue, _ := dataStore.StringForKey(key)
		dataStore.SetString(key, value)

		writer.AddBulkString(existingValue)
	} else {
		writer.AddErrorString("wrong number of arguments for 'set' command")
	}
	writer.WriteResponse()
}
