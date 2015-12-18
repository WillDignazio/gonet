package main

import (
	"digitalbebop.net/gonet/ip"
	"fmt"
	"syscall"
)

func main() {
	sfd, err := ip.OpenRawIPv4Socket()
	if err != nil {
		fmt.Println(err)
		return
	}

	var buff []byte = make([]byte, 1000)
	var buff2 []byte = make([]byte, 1000)

	read, _, _, _, err := syscall.Recvmsg(sfd, buff, buff2, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Read %d bytes", read)
}
