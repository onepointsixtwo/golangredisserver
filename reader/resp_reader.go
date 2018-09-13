package reader

import (
	"io"
)

type RespCommand struct {
	Command string
	Args    []string
}

func CreateRespCommandReader(reader io.Reader) func() (*RespCommand, error) {

	readLine := CreateLineReader(reader)

	return func() (*RespCommand, error) {
		/*
			A RESP COMMAND:

			*3
			$3
			SET
			$5
			mykey
			$7
			myvalue

			Assume each newline represents \r\n (handled by line reader anyway)

			Always starts with a header of *NUM where num is the number of 'parts' afterwards (note distinction between parts and lines here!)

			Each part then has to be parsed as an arg
		*/
		return nil, nil
	}
}
