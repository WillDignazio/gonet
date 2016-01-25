package gonet

import (
	"io"
	"net"
)

type GoNetGatewayType uint

const (
	Uninitialized GoNetGatewayType = 0
	RawSocket     = iota
	Virtual       = iota
)

type GoNetGateway interface {
	io.Reader
	io.Writer
	Name() string
	Interface() net.Interface
	GoNetGatewayType() GoNetGatewayType
}

func GoNetGatewayTypeAsString(t GoNetGatewayType) string {
	switch t {
	case Uninitialized:
		return "Uninitialized"
	case RawSocket:
		return "Raw"
	case Virtual:
		return "Virtual"
	default:
		return "Unknown"
	}
}

