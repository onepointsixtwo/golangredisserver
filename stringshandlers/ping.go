package stringshandlers

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
)

const (
	PONG = "PONG"
)

func pingHandler(args []string, connection connection.Connection) {
	writer := connection.CreateResponseWriter()
	if len(args) > 0 {
		writer.AddBulkString(args[0])
	} else {
		writer.AddSimpleString(PONG)
	}
	writer.WriteResponse()
}
