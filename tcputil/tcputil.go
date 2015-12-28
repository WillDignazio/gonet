package main

import (
	"digitalbebop.net/gonet"
	"fmt"
)

func main() {
	iface, err := gonet.OpenRawSocketInterface()
	if err != nil {
		fmt.Println(err)
		return
	}

	t := gonet.NewGoNetConfig()
	fmt.Println(t)
	fmt.Println(t.GetInt("foo", 10))

	gonet.TestListen(*iface)
}
