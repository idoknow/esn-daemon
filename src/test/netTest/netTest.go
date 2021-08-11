package main

import (
	"esnd/src/service"
	"net"
)

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	service.WriteInt(119812525, c)
}
