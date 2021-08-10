package main

import (
	"esnd/src/service"
	"net"
)

func main() {
	var pa service.PackLogin
	pa.User = "user0"
	pa.Pass = "changeMe"
	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	service.WritePackage(c, pa, 1)

	var p1 service.PackRequest
	p1.From = 0
	p1.Limit = 2

	service.WritePackage(c, p1, 4)
}
