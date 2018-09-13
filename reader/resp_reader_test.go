package reader

import (
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
)

func TestRespReaderCanReadSetCommandCorrectly(t *testing.T) {
	reader := mocks.NewMockReader("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue")
	respRead := CreateRespCommandReader(reader)

	command, err := respRead()

	if err != nil || command == nil {
		t.Error("Command should not be nil and no error should occur parsing first RESP command")
		return
	}

	commandString := command.Command
	firstArg := command.Args[0]
	secondArg := command.Args[1]
	if commandString != "SET" || firstArg != "mykey" || secondArg != "myvalue" {
		t.Errorf("Expected command to be SET with args mykey, myvalue but was %v with args %v, %v", commandString, firstArg, secondArg)
	}
}
