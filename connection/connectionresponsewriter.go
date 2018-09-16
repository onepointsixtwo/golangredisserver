package connection

type ConnectionResponseWriter interface {
	AddSimpleString(str string)
	AddInt(num int)
	AddErrorString(str string)
	AddBulkString(str string)
	ForceArrayType()
	WriteResponse() error
}
