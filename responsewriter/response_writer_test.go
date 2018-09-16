package responsewriter

import (
	"github.com/onepointsixtwo/golangredisserver/connection"
	"testing"
)

// Tests

func TestWritingSimpleStringResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddSimpleString("PONG")
	sut.WriteResponse()

	if mockResponder.responseWritten != "+PONG\r\n" {
		t.Errorf("Expected written response to be +PONG\r\n but was %v", mockResponder.responseWritten)
	}
}

func TestWritingIntResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddInt(2)
	sut.WriteResponse()

	expected := ":2\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written int response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingErrorResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddErrorString("err")
	sut.WriteResponse()

	expected := "-err\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written error response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingNullBulkStringResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddBulkString("")
	sut.WriteResponse()

	expected := "$0\r\n\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written bulk string response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingBulkStringResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddBulkString("this is a bulk string")
	sut.WriteResponse()

	expected := "$21\r\nthis is a bulk string\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written bulk string response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingMixedArrayResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddBulkString("this is a bulk string")
	sut.AddInt(2)
	sut.AddSimpleString("PONG")
	sut.AddErrorString("err")
	sut.WriteResponse()

	expected := "*4\r\n$21\r\nthis is a bulk string\r\n:2\r\n+PONG\r\n-err\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written array response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingEmptyArrayResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.ForceArrayType()
	sut.WriteResponse()

	expected := "*0\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written array response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

func TestWritingSingleElementArrayResponse(t *testing.T) {
	sut, mockResponder := getSut()

	sut.AddInt(2)
	sut.ForceArrayType()
	sut.WriteResponse()

	expected := "*1\r\n:2\r\n"
	if mockResponder.responseWritten != expected {
		t.Errorf("Expected written array response to be %v but was %v", expected, mockResponder.responseWritten)
	}
}

// Helpers
func getSut() (*ResponseWriter, *mockConnection) {
	responder := newMockConnection()
	return New(responder), responder
}

//TODO: move into mocks package - repeated both here and in connectionstore_test.go
// Mock connection
type mockConnection struct {
	responseWritten string
}

func newMockConnection() *mockConnection {
	return &mockConnection{}
}

func (connection *mockConnection) Start() {}

func (connection *mockConnection) SendResponse(response string) {
	connection.responseWritten = response
}

func (connection *mockConnection) CreateResponseWriter() connection.ConnectionResponseWriter {
	return nil
}
