package tcp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// IP Header:
// 	Version: 4
//	IHL: 5 (32bit words)
//	TOS: Precedence(Routine)|Low Delay|Normal Throughput|Normal Reliability
//	Total Length: 52
//	Identification: 63040
//	Flags: Don't Fragment
//	Fragment Offset: 0
//	TTL: 64
//	Protocol: 6
//	Checksum: 64971
//	Source Address: 172.16.119.1
//	Destination Address: 172.16.119.133
var TEST_PACKET1 = []byte{
	0x45, 0x10, 0x00, 0x34, 0xF6, 0x40, 0x40, 0x00,
	0x40, 0x06, 0xFD, 0xCB, 0xAC, 0x10, 0x77, 0x01,
	0xAC, 0x10, 0x77, 0x85, 0xF5, 0x13, 0x00, 0x16,
	0xD6, 0x63, 0x97, 0x05, 0xC3, 0x3A, 0xF8, 0x03,
	0x80, 0x10, 0x0F, 0xFE, 0xB7, 0x0D, 0x00, 0x00,
	0x01, 0x01, 0x08, 0x0A, 0x32, 0xAB, 0x15, 0x66,
	0x00, 0xDA, 0x02, 0x4D,
}

func TestParseIPv4Packet(t *testing.T) {
	packet, err := parseIPv4Message(TEST_PACKET1)
	if err != nil {
		t.Fatal(err.Error())
	}

	onesZero := ^uint16(0)
	fmt.Println("Complement 0: %d", onesZero)

	carry := uint32(0x0A00) + uint32(onesZero)
	fmt.Println("Carry before carry: %d", carry)

	carrybit := ((carry & 0x10000) >> 16)
	fmt.Println("Carry value: %d", carrybit)

	carriedSum := ^uint16(carry & 0xFFFF) + uint16((carry & 0x10000) >> 16)
	fmt.Println("Carried Ones Sum: %d", carriedSum)
	

	fmt.Println("Checksum: %d", CalculateHeaderChecksum(packet))

	assert.Equal(t, uint8(4), packet.Version(), "Version Incorrect")
	assert.Equal(t, uint8(5), packet.IHL(), "Header Length Incorrect")
	assert.Equal(t, uint8(LowDelay), packet.TypeOfService(), "TypeOfService Incorrect")
	assert.Equal(t, uint16(52), packet.TotalLength(), "TotalLength of packet Incorrect")
	assert.Equal(t, uint16(63040), packet.Identification(), "Identification Sequence # Incorrect")
	assert.Equal(t, uint8(DontFragment), packet.Flags(), "Packet Flags Incorrect")
	assert.Equal(t, uint16(16384), packet.FragmentOffset(), "Fragment Offset Incorrect")
	assert.Equal(t, uint8(64), packet.TimeToLive(), "Time to Live Incorrect")
	assert.Equal(t, uint8(6), packet.Protocol(), "IP Protocol incorrect")
	assert.Equal(t, uint16(0xFDCB), packet.HeaderChecksum(), "Header Checksum Incorrect")
}
