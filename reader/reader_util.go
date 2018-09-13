package reader

import (
	"io"
)

const (
	CRLF = "\r\n"
)

func CreateLineReader(reader io.Reader) func() (string, error) {

	// Create the buffer and the line ending we're looking for for command ends
	lineBuffer := make([]string, 0)
	lineEndingBytes := []byte(CRLF)
	buffer := make([]byte, 0)

	return func() (string, error) {
		// Only read from reader when our line buffer is empty.
		for len(lineBuffer) == 0 {
			// Read the next 1024 bytes, and give back any error that occurs
			readBuffer := make([]byte, 1024)
			bytesRead, err := reader.Read(readBuffer)
			if err != nil {
				return "", err
			}

			// Copy bytes across to the buffer, and if we find a line ending, skip the next
			// byte (i.e. \n) and handle the command that's incoming. Also remake the buffer
			// awaiting the next command's ending.
			for i := 0; i < bytesRead; i++ {
				b := readBuffer[i]
				if b == lineEndingBytes[0] {
					command := string(buffer)
					lineBuffer = append(lineBuffer, command)
					buffer = make([]byte, 0)
					i++
				} else {
					buffer = append(buffer, b)
				}
			}
		}

		// When we have a line buffer, pop off the front of the line buffer until it's empty
		line := lineBuffer[0]
		lineBuffer = lineBuffer[1:]

		return line, nil
	}
}
