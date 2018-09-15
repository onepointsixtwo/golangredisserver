package connection

// Connection interface - will later get additions for things like retrieving a name from a connection
// so that the client can find out its connection name
type Connection interface {
	SendResponse(response string)
	CreateResponseWriter() ConnectionResponseWriter
}

type ConnectionResponseWriter interface {
	AddSimpleString(str string)
	AddInt(num int)
	AddErrorString(str string)
	AddBulkString(str string)
	ForceArrayType()
	WriteResponse() error
}
