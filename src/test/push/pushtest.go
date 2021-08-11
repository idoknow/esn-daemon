package main

import (
	"encoding/json"
	"esnd/src/service"
	"fmt"
	"net"
	"time"
)

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}
	var p0 service.PackLogin
	p0.User = "rockchin"
	p0.Pass = "000112rock.,."
	_, err = service.WritePackage(c, p0, 1, "")
	if err != nil {
		panic(err)
	}

	var p1 service.PackNotification
	p1.Target = "soulter,rockchin,root"
	p1.Content = "TestMessage"
	p1.Time = time.Now().String()
	p1.Title = "test"
	_, err = service.WritePackage(c, p1, 3, "")
	if err != nil {
		panic(err)
	}

	var p3 service.PackRequest
	p3.From = 0
	p3.Limit = 100

	_, err = service.WritePackage(c, p3, 4, "")
	if err != nil {
		panic(err)
	}
	for {
		p2json, err := service.ReadPackage(c, "")
		if err != nil {
			panic(err)
		}
		p2 := &service.PackRespNotification{}
		err = json.Unmarshal([]byte(p2json.Json), &p2)
		if err != nil {
			panic(err)
		}

		fmt.Println(" target:" + p2.Target)
	}
}
