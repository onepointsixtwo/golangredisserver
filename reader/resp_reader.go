package reader

import (
	"io"
	"strconv"
	"strings"
)

type RespCommand struct {
	Command string
	Args    []string
}

func CreateRespCommandReader(reader io.Reader) func() (*RespCommand, error) {

	readLine := CreateLineReader(reader)

	return func() (*RespCommand, error) {

		// Get how many 'parts' there are to parse
		partsString, err := readLine()
		if err != nil {
			return nil, err
		}

		parts, err2 := strconv.Atoi(strings.Replace(partsString, "*", "", 1))
		if err2 != nil {
			return nil, err2
		}

		// Parse the command string
		command, err3 := readPart(readLine)
		if err != nil {
			return nil, err3
		}

		// Parse the args
		args := make([]string, 0)
		for i := 0; i < (parts - 1); i++ {
			arg, err4 := readPart(readLine)
			if err4 != nil {
				return nil, err4
			}
			args = append(args, arg)
		}

		return &RespCommand{Command: command, Args: args}, nil
	}
}

func readPart(readLine func() (string, error)) (string, error) {
	line, err := readLine()
	if err != nil {
		return "", err
	}

	strLen, err2 := strconv.Atoi(strings.Replace(line, "$", "", 1))
	if err2 != nil {
		return "", err2
	}

	valueLine, err3 := readLine()
	if err3 != nil {
		return "", err3
	}

	return valueLine[:strLen], nil
}
