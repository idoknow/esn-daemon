package main

import (
	"esnd/src/service"
	"net"
	"sync"
)

func main() {
	var pa service.PackLogin
	pa.User = "user0"
	pa.Pass = "changeMe"
	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	service.WritePackage(c, pa, 1, "")

	var p1 service.PackAccountOperation
	p1.Oper = "add"
	p1.Name = "rockchin"
	p1.Pass = "000112rock.,."
	p1.Priv = "account push pull"
	p1.Kick = true

	service.WritePackage(c, p1, 7, "")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
