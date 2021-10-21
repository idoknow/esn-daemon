package main

import (
	"esnd/src/service"
	"net"
)

func main() {
	var p0 service.PackTest
	p0.Integer = 123456
	p0.Msg = "ThisIsATest MESSAGE!"
	conn, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	service.WriteInt(119812525, conn)
	_, err = service.WritePackage(conn, p0, 0, "")
	if err != nil {
		panic(err)
	}
}
