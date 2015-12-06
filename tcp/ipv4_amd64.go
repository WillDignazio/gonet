package tcp

import (
	"encoding/binary"
)

func (h *IPv4Packet) Version() uint8 {
	return (h.header[0] >> 4) & 0xF
}

func (h *IPv4Packet) IHL() uint8 {
	return h.header[0] & 0xF
}

func (h *IPv4Packet) TypeOfService() uint8 {
	return h.header[1]
}

func (h *IPv4Packet) TotalLength() uint16 {
	return binary.BigEndian.Uint16(h.header[2:4])
}

func (h *IPv4Packet) Identification() uint16 {
	return binary.BigEndian.Uint16(h.header[4:6])
}

func (h *IPv4Packet) Flags() uint8 {
	return h.header[6] >> 5
}

func (h *IPv4Packet) FragmentOffset() uint16 {
	buf := h.header[6:8]
	buf[0] = buf[0] & 0x1F
	return binary.BigEndian.Uint16(buf)
}

func (h *IPv4Packet) TimeToLive() uint8 {
	return h.header[8]
}

func (h *IPv4Packet) Protocol() uint8 {
	return h.header[9]
}

func (h *IPv4Packet) HeaderChecksum() uint16 {
	return binary.BigEndian.Uint16(h.header[10:12])
}

func (h *IPv4Packet) SourceAddress() uint32 {
	return binary.BigEndian.Uint32(h.header[12:16])
}

func (h *IPv4Packet) DestinationAddress() uint32 {
	return binary.BigEndian.Uint32(h.header[16:20])
}

// XXX Preserve byte order so we can use generic flags
func (h *IPv4Packet) Options() uint32 {
	return binary.LittleEndian.Uint32(h.header[20:24]) >> 8
}


func addOnesComplement(x uint16, y uint16) uint16 {
	z := uint32(x) + uint32(y)
	return uint16(z & 0xFFFF) + uint16((z & 0x10000) >> 16)
}

func onesComplementSum(data []byte, base uint16) uint16 {
	var sum uint16 = base
	for idx := 0; idx < len(data); idx += 2 {
		sum = addOnesComplement(sum, (uint16(data[idx]) << 8) | (uint16(data[idx+1])))
	}
	return sum
}

func CalculateHeaderChecksum(packet *IPv4Packet) uint16 {
	var hlen uint8 = packet.IHL() * 8 // 32 bit words / 4 bytes per word
	var checksum uint16

	checksum = onesComplementSum(packet.header[:10], 0)
	checksum = onesComplementSum(packet.header[12:hlen], checksum)

	return ^checksum
}
