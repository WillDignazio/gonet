package gonet

import (
	"fmt"
	"github.com/WillDignazio/gonet/ip"
	"net"
	"os"
	"sync"
	"syscall"
)

type RawSocketInterface struct {
	GoNetGateway
	sockout       int
	sockets       []int
	baseInterface net.Interface
}

func (iface *RawSocketInterface) GoNetGatewayType() GoNetGatewayType {
	return RawSocket
}

func (iface *RawSocketInterface) Interface() net.Interface {
	return iface.baseInterface
}

var DEFAULT_INTERNET_PROTOS = []int{
	syscall.AF_INET,
}

var DEFAULT_TRANSPORT_PROTOS = []int{
	syscall.IPPROTO_TCP,
	syscall.IPPROTO_UDP,
	syscall.IPPROTO_ICMP,
}

func protoToString(proto int) string {
	switch proto {
	case syscall.IPPROTO_TCP:
		return "TCP"
	case syscall.IPPROTO_UDP:
		return "UDP"
	case syscall.IPPROTO_ICMP:
		return "ICMP"
	default:
		return "UNKOWN"
	}
}

func closeRawSocket(sfd int) {
	if err := syscall.Close(sfd); err != nil {
		fmt.Fprint(os.Stderr, "Failed to properly close socket {%d}\n", sfd)
	}
}

func (iface *RawSocketInterface) Read(buffer []byte) (int, error) {
	return 0, nil
}

func (iface *RawSocketInterface) Write(buffer []byte) (int, error) {
	return 0, nil
}

func openRawSocket(net int, proto int) (int, error) {
	fd, err := syscall.Socket(net, syscall.SOCK_RAW, proto)
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

func rawSocketListener(sd int, errch chan<- error) (<-chan []byte, error) {
	var ch chan []byte
	var dummy []byte

	ch = make(chan []byte, 100) // TODO: Tunable
	go func() {
		for {
			buffer := make([]byte, 1500) // TODO: Tunable
			_, _, _, _, err := syscall.Recvmsg(sd, buffer, dummy, syscall.MSG_WAITALL)
			if err != nil {
				errch <- err
				continue
			}

			ch <- buffer
		}
	}()

	return ch, nil
}

func mergeChannels(cs ...<-chan []byte) <-chan []byte {
	var wg sync.WaitGroup
	out := make(chan []byte)
	
	output := func(c <-chan []byte) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func TestListen(iface RawSocketInterface) {
	errs := make(chan error)

	channels := make([]<-chan []byte, len(iface.sockets))
	for idx, val := range iface.sockets {
		ch, err := rawSocketListener(val, errs)
		if err != nil {
			fmt.Println("BARRRRR", err)
			return
		}
		fmt.Println("Created ", ch)
		channels[idx] = ch
	}

	ch := mergeChannels(channels...)

	for {
		_, err := ip.ParseIPv4Datagram(<-ch)
		if err != nil {
			errs <- err
			continue
		}
	}
}

func OpenRawSocketInterface() (*RawSocketInterface, error) {
	var sockets []int

	sockout, err := openRawSocket(syscall.AF_INET, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, err
	}

	sockets = make([]int, len(DEFAULT_TRANSPORT_PROTOS)*len(DEFAULT_INTERNET_PROTOS))
	for pdx, ipProto := range DEFAULT_INTERNET_PROTOS {
		for tdx, transProto := range DEFAULT_TRANSPORT_PROTOS {
			sock, err := openRawSocket(ipProto, transProto)
			if err != nil {
				for sidx := 0; sidx < tdx; sidx += 1 {
					closeRawSocket(sockets[sidx])
				}
				return nil, err
			}
			sockets[pdx*len(DEFAULT_TRANSPORT_PROTOS)+tdx] = sock
		}
	}

	mac, _ := net.ParseMAC("01:23:45:67:89:ab")
	return &RawSocketInterface{
		sockout: sockout,
		sockets: sockets,
		baseInterface: net.Interface{ // XXX TODO: Replace with real interface
			Index:        0,
			MTU:          (15 * _KB_),
			Name:         "gn0",
			HardwareAddr: mac,
			Flags:        net.FlagUp,
		},
	}, nil
}
