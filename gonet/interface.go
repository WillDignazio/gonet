package gonet

import (
	"io"
	"net"
)

type GoNetInterfaceType uint

const (
	Uninitialized GoNetInterfaceType = 0
	RawSocket     = iota
	Virtual       = iota
)

type GoNetInterface interface {
	io.Reader
	io.Writer
	Name() string
	Interface() net.Interface
	GoNetInterfaceType() GoNetInterfaceType
}

func GoNetInterfaceTypeAsString(t GoNetInterfaceType) string {
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

