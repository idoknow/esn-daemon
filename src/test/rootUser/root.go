package main

import (
	"esnd/src/service"
	"net"
	"sync"
)

func main() {
	var pa service.PackLogin
	pa.User = "root"
	pa.Pass = "changeMe"
	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	service.WriteInt(119812525, c)
	service.ReadInt(c)
	service.WritePackage(c, pa, 1, "")

	var p1 service.PackAccountOperation
	p1.Oper = "remove"
	p1.Name = "fuckyou"
	p1.Pass = "000112rock.,."
	p1.Priv = "account push pull"
	p1.Kick = true

	service.WritePackage(c, p1, 7, "")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
