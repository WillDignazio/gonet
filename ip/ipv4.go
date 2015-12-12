package tcp

// Architecture dependent getters are provided in ipv4_<arch>, we only grab the
// data stream in this file so we don't do uneccessary byte conversions or data
// allocations.

import (
	"errors"
	"fmt"
	"net"
	"syscall"
)

type IPv4Header interface {
	Version() uint8
	IHL() uint8
	TypeOfService() uint8 // XXX Bitfield to remain in network order

	// TOS -> DSCP in RFC 2474, support for DiffServ
	// (Differentiated Services). We're going to leave the original
	// TypeOfService interface intact, and as they're shared bit positions.
	DifferentiatedServicesCodePoint() uint8

	// TOS -> ECN defined in RFC 3168, allows end-to-end notification of
	// network congestion without dropping packets.
	ExplicitCongestionNotification() uint8

	TotalLength() uint16
	Identification() uint16
	Flags() uint8
	FragmentOffset() uint16
	TimeToLive() uint8
	Protocol() uint8
	HeaderChecksum() uint16
	SourceAddress() net.IP
	DestinationAddress() net.IP
	Options() []byte

	CalculateChecksum() uint16
}

type IPv4Packet struct {
	header  []byte
	data    []byte
}

var _ IPv4Header = (*IPv4Packet)(nil) // Enforce that we have an impl

// Minimum size of the header
const IPV4_HEADER_PREAMBLE_SIZE = 20

// Bits 0-2:  Precedence.
// Bit    3:  0 = Normal Delay,      1 = Low Delay.
// Bits   4:  0 = Normal Throughput, 1 = High Throughput.
// Bits   5:  0 = Normal Reliability, 1 = High Relibility.
// Bit  6-7:  Reserved for Future Use.
//
//    0     1     2     3     4     5     6     7
// +-----+-----+-----+-----+-----+-----+-----+-----+
// |                 |     |     |     |     |     |
// |   PRECEDENCE    |  D  |  T  |  R  |  0  |  0  |
// |                 |     |     |     |     |     |
// +-----+-----+-----+-----+-----+-----+-----+-----+
//
//   Precedence
//
//     111 - Network Control
//     110 - Internetwork Control
//     101 - CRITIC/ECP
//     100 - Flash Override
//     011 - Flash
//     010 - Immediate
//     001 - Priority
//     000 - Routine
type IPv4PrecedenceMask uint8

const (
	NetworkControl      IPv4PrecedenceMask = 0xE0
	InternetworkControl IPv4PrecedenceMask = 0xC0
	CRITIC_ECP          IPv4PrecedenceMask = 0xA0
	FlashOverride       IPv4PrecedenceMask = 0x80
	Flash               IPv4PrecedenceMask = 0x60
	Immediate           IPv4PrecedenceMask = 0x40
	Priority            IPv4PrecedenceMask = 0x20
	Routine             IPv4PrecedenceMask = 0x00
)

type IPv4ServiceTypeMask uint8

const (
	HighReliability IPv4ServiceTypeMask = 1 << 2
	HighThroughput  IPv4ServiceTypeMask = 1 << 3
	LowDelay        IPv4ServiceTypeMask = 1 << 4
)

// Flags:  3 bits
//
//   Various Control Flags.
//
//     Bit 0: reserved, must be zero
//     Bit 1: (DF) 0 = May Fragment,  1 = Don't Fragment.
//     Bit 2: (MF) 0 = Last Fragment, 1 = More Fragments.
//
//         0   1   2
//       +---+---+---+
//       |   | D | M |
//       | 0 | F | F |
//       +---+---+---+
type IPv4Flag uint8

const (
	MoreFragments IPv4Flag = 1 << 0
	DontFragment  IPv4Flag = 1 << 1
)

func (h *IPv4Packet) Options() []byte {
	ihl := h.IHL()
	return h.header[IPV4_HEADER_PREAMBLE_SIZE:(ihl*4)]
}

func (h *IPv4Packet) SourceAddress() net.IP {
	return h.header[12:16]
}

func (h *IPv4Packet) DestinationAddress() net.IP {
	return h.header[16:20]
}


// Don't want to waste time building and tearing down the packet objects
// we're just going to extract it outright
func extractIHL(data []byte) uint8 {
	return data[0] & 0x0F
}

func parseIPv4Message(rawData []byte) (*IPv4Packet, error) {
	rawDatalen := len(rawData)
	if rawDatalen < IPV4_HEADER_PREAMBLE_SIZE {
		msg := fmt.Sprintf("Too small buffer size, can't parse message: %d of %d",
			rawDatalen, IPV4_HEADER_PREAMBLE_SIZE)
		return nil, errors.New(msg)
	}

	ihl := extractIHL(rawData)
	eoh := ihl * 4
	packet := IPv4Packet{
		header: rawData[0:eoh],
		data: rawData[eoh:],
	}

	return &packet, nil
}

func OpenRawIPv4Socket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return -1, err
	}

	// Configure the socket so that we must manually include the IP headers
	err = syscall.SetsockoptInt(fd, 0, syscall.IP_HDRINCL, 1)
	if err != nil {
		return -1, err
	}

	return fd, nil
}

