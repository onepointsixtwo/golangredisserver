package connection

type Connection interface {
	Start()
	SendResponse(response string)
	CreateResponseWriter() ConnectionResponseWriter
}
