package tcp

import (
	"fmt"
	"testing"
)

func TestParseSegment(t *testing.T) {
	test := TCPSegment{SequenceNumber: 0}
	fmt.Println("test: %d", test.SequenceNumber)
}
