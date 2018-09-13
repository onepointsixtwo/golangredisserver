package mocks

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// Tests
func TestMockConnectionCanBeTestToNetConn(t *testing.T) {
	var conn net.Conn
	conn = NewMockConnection("PING\r\n")

	if conn == nil {
		t.Error("MockConnection should be castable to net.Conn")
	}
}

func TestMockConnectionUsesReadStringForReadContent(t *testing.T) {
	inputCommands := "PING\r\n"
	conn := NewMockConnection(inputCommands)
	outputCommands := readAllConnectionInputToString(conn)
	if inputCommands != outputCommands {
		t.Errorf("Expected read output to be %v but was %v", inputCommands, outputCommands)
	}
}

func TestWriteInputsDataToWriteBuffer(t *testing.T) {
	outputString := "+PONG"
	bytes := []byte(outputString)

	conn := NewMockConnection("PING\r\n")
	conn.Write(bytes)

	outputStringWritten := conn.WriteBuffer.String()

	if outputStringWritten != outputString {
		t.Errorf("Expected an output string of %v to be written to buffer, but was %v", outputString, outputStringWritten)
	}
}

func TestSettingDeadlines(t *testing.T) {
	deadline := time.Now()

	conn := NewMockConnection("")
	conn.SetDeadline(deadline)

	if conn.ReadDeadline != deadline || conn.WriteDeadline != deadline {
		t.Error("Read or write deadline not set correctly")
	}
}

func TestClosing(t *testing.T) {
	conn := NewMockConnection("")
	conn.Close()

	if conn.Closed != true {
		t.Error("Connection should have closed value of true after conn.Close() called")
	}
}

// Helpers

func readAllConnectionInputToString(conn net.Conn) string {
	var buffer bytes.Buffer

	for {
		// Read the next 1024 bytes and break on errors.
		readBuffer := make([]byte, 1024)
		bytesRead, err := conn.Read(readBuffer)
		if err != nil {
			break
		}

		for i := 0; i < bytesRead; i++ {
			buffer.WriteByte(readBuffer[i])
		}
	}

	return buffer.String()
}
