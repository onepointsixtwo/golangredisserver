package stringshandlers

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"github.com/onepointsixtwo/golangredisserver/expiry"
	"strconv"
)

func expireHandler(args []string, connection connection.Connection, expiryHandler *expiry.Handler) {
	writer := connection.CreateResponseWriter()
	if len(args) == 2 {
		key := args[0]
		expirySecondsString := args[1]

		expirySeconds, err := strconv.Atoi(expirySecondsString)
		if err != nil {
			writer.AddErrorString(fmt.Sprintf("unable to parse argument for expiry time in seconds: %v", err))
		} else {
			err := expiryHandler.ExpireKeyAfterSeconds(key, expirySeconds)
			if err != nil {
				writer.AddErrorString("cannot set expiry for non existent key!")
			} else {

				writer.AddSimpleString(OK)
			}
		}
	} else {
		writer.AddErrorString("incorrect number of args - should have two, a key and expiry time")
	}
	writer.WriteResponse()
}
