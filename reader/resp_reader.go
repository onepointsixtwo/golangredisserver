package reader

import (
	"fmt"
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
		parts, err := getPartsCount(readLine)
		if err != nil {
			return nil, err
		}

		// Parse the args
		commandAndArgs := make([]string, 0)
		for i := 0; i < parts; i++ {
			part, readPartErr := readPart(readLine)
			if readPartErr != nil {
				return nil, readPartErr
			}
			commandAndArgs = append(commandAndArgs, part)
		}

		if len(commandAndArgs) > 0 {
			return &RespCommand{Command: commandAndArgs[0], Args: commandAndArgs[1:]}, nil
		} else {
			return nil, fmt.Errorf("Command and args have zero members - command cannot be read! Should be %v parts", parts)
		}
	}
}

func getPartsCount(readLine func() (string, error)) (int, error) {
	// Get how many 'parts' there are to parse
	partsString, err := readLine()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.Replace(partsString, "*", "", 1))
}

func readPart(readLine func() (string, error)) (string, error) {
	line, err := readLine()
	if err != nil {
		return "", err
	}

	// If the prefix is '$' it's a bulk string split over 2 lines, but if it's another then it's all on one line.
	if strings.HasPrefix(line, "$") {
		strLen, err2 := strconv.Atoi(strings.Replace(line, "$", "", 1))
		if err2 != nil {
			return "", err2
		}

		valueLine, err3 := readLine()
		if err3 != nil {
			return "", err3
		}

		return valueLine[:strLen], nil
	} else if strings.HasPrefix(line, ":") || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
		return line[1:], nil
	}

	return "", fmt.Errorf("Error: unknown prefix for part %v", line)
}
