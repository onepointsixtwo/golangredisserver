package stringshandlers

import (
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/connection"
	"time"
)

func timeHandler(args []string, connection connection.Connection) {
	writer := connection.CreateResponseWriter()

	currentTime := time.Now()

	//Get the seconds
	seconds := currentTime.Unix()
	writer.AddBulkString(fmt.Sprintf("%v", seconds))

	//Get the microseconds
	nanoseconds := currentTime.UnixNano()
	nanosecondsRemainder := nanoseconds % (seconds * int64(time.Nanosecond))
	milliseconds := nanosecondsRemainder / 1000
	writer.AddBulkString(fmt.Sprintf("%v", milliseconds))

	writer.WriteResponse()
}
