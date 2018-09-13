package mocks

import (
	"io"
)

type MockReader struct {
	value  []byte
	offset int
}

func NewMockReader(value string) *MockReader {
	return &MockReader{value: []byte(value)}
}

func (reader *MockReader) Read(b []byte) (n int, err error) {
	lengthOfGivenArray := len(b)

	totalReadStringLength := len(reader.value)
	remainingReadString := totalReadStringLength - reader.offset

	if remainingReadString <= 0 {
		return 0, io.EOF
	} else {
		lengthToRead := remainingReadString
		if lengthOfGivenArray < lengthToRead {
			lengthToRead = lengthOfGivenArray
		}

		position := 0
		for i := reader.offset; i < (reader.offset + lengthToRead); i++ {
			b[position] = reader.value[i]
			position++
		}

		reader.offset += lengthToRead

		return lengthToRead, nil
	}
}
