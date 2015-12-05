package tcp

// Architecture dependent getters are provided in ipv4_<arch>, we only grab the
// data stream in this file so we don't do uneccessary byte conversions or data
// allocations.

import (
	"errors"
	"fmt"
	"syscall"
)

const IPV4_BLOCK_SIZE = 64 * 1024
const IPV4_HEADER_LENGTH = 24 // 24 Octets

type IPv4Precedence uint8
const (
	NetworkControl IPv4Precedence = 0x70
	InternetworkControl IPv4Precedence = 0x60
	CRITIC_ECP IPv4Precedence = 0x50
	FlashOverride IPv4Precedence = 0x40
	Flash IPv4Precedence = 0x30
	Immediate IPv4Precedence = 0x20
	Priority = 0x10
	Routine = 0
)

type IPv4ServiceType uint8
const (
	LowDelay IPv4ServiceType = 1 << 3
	HighThroughput IPv4ServiceType = 1 << 4
	HighReliability IPv4ServiceType = 1 << 5
)

type IPv4Flag uint8
const (

	
	// Flags
	DontFragment  IPv4Flag = 1 << 1
	MoreFragments IPv4Flag = 1 << 2
)

type IPv4Packet struct {
	header []byte
	data   []byte
}

func parseIPv4Message(rawData []byte) (*IPv4Packet, error) {
	rawDatalen := len(rawData)
	if rawDatalen < IPV4_HEADER_LENGTH {
		return nil, errors.New(fmt.Sprintf("Too small buffer size, can't parse message: %d of %d", rawDatalen, IPV4_HEADER_LENGTH))
	}

	packet := IPv4Packet{
		header: rawData[0:IPV4_HEADER_LENGTH],
		data:   rawData[IPV4_HEADER_LENGTH:],
	}

	return &packet, nil
}

func OpenRawIPv4Socket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return -1, err
	}

	return fd, nil
}
