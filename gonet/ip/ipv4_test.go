package ip

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// IP Header:
//	Version: 4
//	IHL: 8 (32bit words)
//	TOS: Precedence(Normal)|Normal Delay|Normal Throughput|Normal Reliability
//	Total Length: 576
//	Identification: 111
//	Flags:
//	Fragment Offset: 0
//	TTL: 123
//	Protocol: 6
//	Checksum:
//	Source Address: 172.16.119.1
//	Destination Address: 172.16.119.133
var TEST_PACKET2 = []byte{
	0x48, 0x00, 0x02, 0x40, 0x00, 0x6F, 0x00, 0x00,
	0x7B, 0x06, 0xF3, 0xA1, 0xAC, 0x10, 0x77, 0x01,
	0xAC, 0x10, 0x77, 0x85, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// Example packet from RFC791, page 38
	// Data left out
}

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

	assert.Equal(t, uint8(4), packet.Version(), "Version incorrect")
	assert.Equal(t, uint8(5), packet.IHL(), "Header length incorrect")
	assert.Equal(t, uint8(LowDelay), packet.TypeOfService(), "TypeOfService incorrect")
	assert.Equal(t, uint16(52), packet.TotalLength(), "TotalLength of packet incorrect")
	assert.Equal(t, uint16(63040), packet.Identification(), "Identification sequence # incorrect")
	assert.Equal(t, uint8(DontFragment), packet.Flags(), "Packet flags incorrect")
	assert.Equal(t, uint16(0), packet.FragmentOffset(), "Fragment offset incorrect")
	assert.Equal(t, uint8(64), packet.TimeToLive(), "Time to live incorrect")
	assert.Equal(t, uint8(6), packet.Protocol(), "IP protocol number incorrect")
	assert.Equal(t, uint16(0xFDCB), packet.HeaderChecksum(), "Header checksum incorrect")
	assert.Equal(t, uint16(0xFDCB), packet.CalculateChecksum(), "Calculated checksum incorrect")
	assert.Equal(t, true, packet.Valid(), "Invalid checksum found for packet")
	assert.Equal(t, "172.16.119.1", packet.SourceAddress().String(), "Invalid source address")
	assert.Equal(t, "172.16.119.133", packet.DestinationAddress().String(), "Invalid destination address")
}

func TestParseIPv4PacketWithOptions(t *testing.T) {
	packet, err := parseIPv4Message(TEST_PACKET2)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, uint8(4), packet.Version(), "Version incorrect")
	assert.Equal(t, uint8(8), packet.IHL(), "Header length incorrect")
	assert.Equal(t, uint8(0), packet.TypeOfService(), "TypeOfService incorrect")
	assert.Equal(t, uint16(576), packet.TotalLength(), "TotalLength of packet incorrect")
	assert.Equal(t, uint16(111), packet.Identification(), "Identification sequence # incorrect")
	assert.Equal(t, uint8(0), packet.Flags(), "Packet flags incorrect")
	assert.Equal(t, uint16(0), packet.FragmentOffset(), "Fragment offset incorrect")
	assert.Equal(t, uint8(123), packet.TimeToLive(), "Time to live incorrect")
	assert.Equal(t, uint8(6), packet.Protocol(), "IP protocol number incorrect")
	assert.Equal(t, uint16(0xF3A1), packet.HeaderChecksum(), "Header checksum incorrect")
	assert.Equal(t, uint16(0xF3A1), packet.CalculateChecksum(), "Calculated checksum incorrect")
	assert.Equal(t, true, packet.Valid(), "Invalid checksum found for packet")
	assert.Equal(t, "172.16.119.1", packet.SourceAddress().String(), "Invalid source address")
	assert.Equal(t, "172.16.119.133", packet.DestinationAddress().String(), "Invalid destination address")
}
