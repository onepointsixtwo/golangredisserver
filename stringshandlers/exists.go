package stringshandlers

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
)

func existsHandler(args []string, connection connection.Connection, dataStore keyvaluestore.Store) {
	writer := connection.CreateResponseWriter()

	exists := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		_, err := dataStore.StringForKey(key)
		if err == nil {
			exists++
		}
	}

	writer.AddInt(exists)
	writer.WriteResponse()
}
