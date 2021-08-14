package main

import (
	"esnd/src/service"
	"net"
	"sync"
)

func main() {
	var pa service.PackLogin
	pa.User = "root"
	pa.Pass = "turtle"
	c, err := net.Dial("tcp", "39.100.5.139:3003")
	if err != nil {
		panic(err)
	}
	service.WriteInt(119812525, c)
	service.ReadInt(c)
	service.WritePackage(c, pa, 1, "")

	var p1 service.PackAccountOperation
	p1.Oper = "add"
	p1.Name = "autoxmrig"
	p1.Pass = "000112rock.,."
	p1.Priv = "push pull"
	p1.Kick = true

	service.WritePackage(c, p1, 7, "")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
