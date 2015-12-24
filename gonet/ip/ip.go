package ip

type IPDatagram interface {
	Header() []byte
	Data() []byte
}

