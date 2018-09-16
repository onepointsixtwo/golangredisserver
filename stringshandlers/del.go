package stringshandlers

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"github.com/onepointsixtwo/golangredisserver/keyvaluestore"
)

func deleteHandler(args []string, connection connection.Connection, dataStore keyvaluestore.Store, expiryHandler *expiry.Handler) {
	writer := connection.CreateResponseWriter()

	deleted := 0
	for i := 0; i < len(args); i++ {
		key := args[i]
		success := dataStore.DeleteString(key)
		if success {
			expiryHandler.CancelTimerForKeyIfExists(key)
			deleted++
		}
	}

	writer.AddInt(deleted)
	writer.WriteResponse()
}
