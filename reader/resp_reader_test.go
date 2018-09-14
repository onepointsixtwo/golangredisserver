package reader

import (
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
)

func TestRespReaderCanReadSetCommandCorrectly(t *testing.T) {
	reader := mocks.NewMockReader("*4\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n:5\r\n")
	respRead := CreateRespCommandReader(reader)

	command, err := respRead()

	if err != nil || command == nil {
		t.Errorf("Command should not be nil and no error should occur parsing first RESP command %v", err)
		return
	}

	commandString := command.Command
	firstArg := command.Args[0]
	secondArg := command.Args[1]
	thirdArg := command.Args[2]
	if commandString != "SET" || firstArg != "mykey" || secondArg != "myvalue" || thirdArg != "5" {
		t.Errorf("Expected command to be SET with args mykey, myvalue but was %v with args %v, %v, %v", commandString, firstArg, secondArg, thirdArg)
	}
}
