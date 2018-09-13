package reader

import (
	"github.com/onepointsixtwo/golangredisserver/mocks"
	"testing"
)

func TestLineReader(t *testing.T) {
	reader := mocks.NewMockReader("Some value to read\r\nSome other line\r\nSome third line\r\n")

	lineRead := CreateLineReader(reader)

	first, _ := lineRead()
	if first != "Some value to read" {
		t.Errorf("Expected first to be 'Some value to read' but was %v", first)
	}

	second, _ := lineRead()
	if second != "Some other line" {
		t.Errorf("Expected second to be 'Some other line' but was %v", second)
	}

	third, _ := lineRead()
	if third != "Some third line" {
		t.Errorf("Expected second to be 'Some third line' but was %v", third)
	}

	_, err := lineRead()
	if err == nil {
		t.Error("Expected error (io.EOF) after all lines read")
	}
}
