package main

import (
	"digitalbebop.net/tcp"
	"fmt"
)

func main() {
	_, err := tcp.OpenRawIPv4Socket()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Hello")
}
