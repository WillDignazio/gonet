package tcp

import (
	"encoding/binary"
)

func (h *IPv4Packet) Version() uint8 {
	return (h.header[0] >> 4) & 0xF
}

// XXX Bitfield to remain in network order
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
