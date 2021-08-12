package main

import (
	"encoding/json"
	"esnd/src/service"
	"fmt"
	"net"
	"sync"
)

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	service.WriteInt(119812525, c)
	service.ReadInt(c)
	var p0 service.PackReqRSAKey
	service.WritePackage(c, p0, 8, "")

	service.ReadPackage(c, "")

	p1, err := service.ReadPackage(c, "")
	if err != nil {
		panic(err)
	}
	p2 := &service.PackRSAPublicKey{}
	err = json.Unmarshal([]byte(p1.Json), &p2)
	if err != nil {
		panic(err)
	}
	fmt.Println("key:" + p2.PublicKey)

	var p3 service.PackLogin
	p3.User = "fuckyou"
	p3.Pass = "fuckyou"
	_, err = service.WritePackage(c, p3, 1, p2.PublicKey)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
