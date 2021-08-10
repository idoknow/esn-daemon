package main

import (
	"esnd/src/db"
	"esnd/src/service"
	"esnd/src/util"
	"strconv"
	"sync"
)

var Cfg *util.Config
var Service *service.NetService
var wg sync.WaitGroup

func main() {
	Cfg, err := util.LoadConfig(".esnd/esnd.conf")
	if err != nil {
		panic(err)
	}
	util.SaySub("Main", "Config file loaded.")

	err = db.Init(Cfg)
	if err != nil {
		panic(err)
	}
	util.SaySub("Main", "Database loaded.")

	port, err := strconv.Atoi(Cfg.GetAnyway("service.port", "3003"))
	if err != nil {
		panic(err)
	}
	Service, err = service.MakeNS(port)
	if err != nil {
		panic(err)
	}
	go Service.Accept()

	wg.Add(1)
	wg.Wait()
}
