package responsewriter

import (
	"bytes"
	"fmt"
	"github.com/onepointsixtwo/golangredisserver/router"
)

const (
	CRLF = "\r\n"
)

// Types
type ResponseWriter struct {
	responder      router.Responder
	responses      []interface{}
	forceArrayType bool
}

type iSimpleString interface {
	valueFromSimpleString() string
}

type simpleString struct {
	value string
}

type iRespInt interface {
	valueFromRespInt() string
}

type respInt struct {
	value int
}

type iErrorString interface {
	valueFromErrorString() string
}

type errorString struct {
	value string
}

type iBulkString interface {
	valueFromBulkString() string
}

type bulkString struct {
	value string
}

// Public methods
func New(responder router.Responder) *ResponseWriter {
	return &ResponseWriter{responder: responder, responses: make([]interface{}, 0)}
}

func (writer *ResponseWriter) AddSimpleString(str string) {
	value := &simpleString{str}
	writer.responses = append(writer.responses, value)
}

func (writer *ResponseWriter) AddInt(num int) {
	value := &respInt{num}
	writer.responses = append(writer.responses, value)
}

func (writer *ResponseWriter) AddErrorString(str string) {
	value := &errorString{str}
	writer.responses = append(writer.responses, value)
}

func (writer *ResponseWriter) AddBulkString(str string) {
	value := &bulkString{str}
	writer.responses = append(writer.responses, value)
}

func (writer *ResponseWriter) ForceArrayType() {
	writer.forceArrayType = true
}

func (writer *ResponseWriter) WriteResponse() error {
	responsesCount := len(writer.responses)

	if responsesCount == 1 && !writer.forceArrayType {
		writer.writeSingleResponse(writer.responses[0])
	} else if responsesCount > 1 || writer.forceArrayType {
		writer.writeArrayResponse(writer.responses)
	} else {
		return fmt.Errorf("Cannot formulate a response with 0 responses %v", writer.responses)
	}

	return nil
}

func (writer *ResponseWriter) writeSingleResponse(value interface{}) {
	writer.responder.SendResponse(writer.stringValueFromInterface(value))
}

func (writer *ResponseWriter) writeArrayResponse(responses []interface{}) {
	var buffer bytes.Buffer

	responsesLength := len(responses)

	fmt.Fprintf(&buffer, "*%v%v", responsesLength, CRLF)

	for i := 0; i < responsesLength; i++ {
		value := responses[i]
		fmt.Fprintf(&buffer, "%v", writer.stringValueFromInterface(value))
	}

	writer.responder.SendResponse(buffer.String())
}

func (writer *ResponseWriter) stringValueFromInterface(value interface{}) string {
	switch v := value.(type) {
	case iSimpleString:
		return v.valueFromSimpleString()
	case iRespInt:
		return v.valueFromRespInt()
	case iErrorString:
		return v.valueFromErrorString()
	case iBulkString:
		return v.valueFromBulkString()
	default:
		fmt.Printf("Returning default - found %v %v \n", value, v)
		return ""
	}
}

func (simple *simpleString) valueFromSimpleString() string {
	return fmt.Sprintf("+%v%v", simple.value, CRLF)
}

func (rInt *respInt) valueFromRespInt() string {
	return fmt.Sprintf(":%v%v", rInt.value, CRLF)
}

func (err *errorString) valueFromErrorString() string {
	return fmt.Sprintf("-%v%v", err.value, CRLF)
}

func (bulk *bulkString) valueFromBulkString() string {
	return fmt.Sprintf("$%v%v%v%v", len(bulk.value), CRLF, bulk.value, CRLF)
}
