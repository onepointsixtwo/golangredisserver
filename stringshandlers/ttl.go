package stringshandlers

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
)

func ttlHandler(args []string, connection connection.Connection, expiryHandler *expiry.Handler) {
	writer := connection.CreateResponseWriter()
	if len(args) == 1 {
		key := args[0]
		ttl, err := expiryHandler.RemainingExpiryTTLForKey(key)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("no expiry time exists for key %v", key))
		} else {
			writer.AddInt(ttl)
		}
	} else {
		writer.AddErrorString(fmt.Sprintf("incorrect number of args for TTL - expected 1 but got %v", len(args)))
	}
	writer.WriteResponse()
}
