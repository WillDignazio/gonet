package tcp

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Generated Using rpi (github.com/WillDignazio/rpi.c)
// IP Header:
// 	Version: 4
//	IHL: 5 (32bit words)
//	TOS: Precedence(Priority)|Normal Delay|Normal Throughput|Normal Reliability
//	Total Length: 52
//	Identification: 63040
//	Flags: Don't Fragment|More Fragments
//	Fragment Offset: 16384
//	TTL: 64
//	Protocol: 6
//	Checksum: 52221
//	Source Address: 172.16.119.1
//	Destination Address: 172.16.119.133
var TEST_PACKET1 = []byte{
	0x45, 0x10, 0x00, 0x34, 0xF6, 0x40, 0x40, 0x00, 0x40,
	0x06, 0xFD, 0xCB, 0xAC, 0x10, 0x77, 0x01, 0xAC, 0x10,
	0x77, 0x85, 0xF5, 0x13, 0x00, 0x16, 0xD6, 0x63, 0x97,
	0x05, 0xC3, 0x3A, 0xF8, 0x03, 0x80, 0x10, 0x0F, 0xFE,
	0xB7, 0x0D, 0x00, 0x00, 0x01, 0x01, 0x08, 0x0A, 0x32,
	0xAB, 0x15, 0x66, 0x00, 0xDA, 0x02, 0x4D,
}


func TestParseIPv4Packet(t *testing.T) {
	packet, err := parseIPv4Message(TEST_PACKET1)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println("TotalLength: %d", packet.TotalLength())
	fmt.Println("version: %d", packet.Version())

	assert.Equal(t, uint8(4), packet.Version(), "Version Incorrect")
	assert.Equal(t, uint8(5), packet.IHL(), "Header Length Incorrect")
	assert.Equal(t, uint8(Priority), packet.TypeOfService(), "TypeOfService Incorrect")
	assert.Equal(t, uint16(52), packet.TotalLength(), "TotalLength of packet Incorrect")
	assert.Equal(t, uint16(63040), packet.Identification(), "Identification Sequence # Incorrect")
	assert.Equal(t, MoreFragments, packet.Flags(), "Incorrect Packet Flags")
}
