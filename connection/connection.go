package connection

type Connection interface {
	SendResponse(response string)
	CreateResponseWriter() ConnectionResponseWriter
}
