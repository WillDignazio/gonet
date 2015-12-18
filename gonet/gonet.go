package net

import (
	"digitalbebop.net/gonet/ip"
	"fmt"
)

func Test() {
	_, e := ip.OpenRawIPv4Socket()
	if e != nil {
		fmt.Println("Fuck")
		return
	}

	fmt.Println(":)")
}
