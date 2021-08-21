package main

import (
	"encoding/json"
	"esnd/src/service"
	"fmt"
	"net"
	"strconv"
)

func main() {

	c, err := net.Dial("tcp", "127.0.0.1:3003")
	if err != nil {
		panic(err)
	}

	err = service.WriteInt(119812525, c)
	if err != nil {
		panic(err)
	}
	fmt.Println("protocol:" + strconv.Itoa(service.ReadInt(c)))

	var p0 service.PackLogin
	p0.User = "root"
	p0.Pass = "changeMe"
	_, err = service.WritePackage(c, p0, 1, "")
	if err != nil {
		panic(err)
	}

	var p1 service.PackReqRecent
	p1.Limit = 10
	p1.Token = "sah 9udhsiuhasoudha"

	service.WritePackage(c, p1, 10, "")

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

		fmt.Println("id:" + strconv.Itoa(p2.Id) + " target:" + p2.Target + " content:" + p2.Content)
	}
}
