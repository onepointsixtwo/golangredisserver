package stringshandlers

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
)

func getHandler(args []string, connection connection.Connection, dataStore keyvaluestore.Store) {
	writer := connection.CreateResponseWriter()
	if len(args) > 0 {
		key := args[0]
		value, err := dataStore.StringForKey(key)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("value not found for key '%v'", key))
		} else {
			writer.AddBulkString(value)
		}
	} else {
		writer.AddErrorString("wrong number of arguments for 'get' command")
	}
	writer.WriteResponse()
}
