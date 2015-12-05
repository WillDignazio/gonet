package tcp

type octet uint8

type TCPSegment struct {
	SequenceNumber octet
}
