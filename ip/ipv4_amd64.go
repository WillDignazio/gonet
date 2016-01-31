package ip

import (
	"encoding/binary"
)

func (h *IPv4Datagram) Version() uint8 {
	return (h.header[0] >> 4) & 0xF
}

func (h *IPv4Datagram) IHL() uint8 {
	return h.header[0] & 0xF
}

func (h *IPv4Datagram) TypeOfService() uint8 {
	return h.header[1]
}

func (h *IPv4Datagram) DifferentiatedServicesCodePoint() uint8 {
	return (h.header[1] >> 2) & 0x3F
}

func (h *IPv4Datagram) ExplicitCongestionNotification() uint8 {
	return h.header[1] & 0x03
}

func (h *IPv4Datagram) TotalLength() uint16 {
	return binary.BigEndian.Uint16(h.header[2:4])
}

func (h *IPv4Datagram) Identification() uint16 {
	return binary.BigEndian.Uint16(h.header[4:6])
}

func (h *IPv4Datagram) Flags() uint8 {
	return h.header[6] >> 5
}

func (h *IPv4Datagram) FragmentOffset() uint16 {
	return binary.BigEndian.Uint16([]byte{h.header[6] & 0x1F, h.header[7]})
}

func (h *IPv4Datagram) TimeToLive() uint8 {
	return h.header[8]
}

func (h *IPv4Datagram) Protocol() uint8 {
	return h.header[9]
}

func (h *IPv4Datagram) HeaderChecksum() uint16 {
	return binary.BigEndian.Uint16(h.header[10:12])
}

func addOnesComplement(x uint16, y uint16) uint16 {
	z := uint32(x) + uint32(y)
	return uint16(z&0xFFFF) + uint16((z&0x10000)>>16)
}

func onesComplementSum(data []byte, base uint16) uint16 {
	var sum uint16 = base
	for idx := 0; idx < len(data); idx += 2 {
		sum = addOnesComplement(sum, (uint16(data[idx])<<8)|(uint16(data[idx+1])))
	}
	return sum
}

func (packet *IPv4Datagram) CalculateChecksum() uint16 {
	var hlen uint8 = packet.IHL() * 4 // 32 bit words / 4 bytes per word
	var checksum uint16
	checksum = onesComplementSum(packet.header[:10], 0)
	checksum = onesComplementSum(packet.header[12:hlen], checksum)
	return ^checksum
}

func (packet *IPv4Datagram) Valid() bool {
	var hlen uint8 = packet.IHL() * 4
	chk := ^onesComplementSum(packet.header[:hlen], 0)
	return chk == 0
}
