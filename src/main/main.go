package main

import (
	"bufio"
	"esnd/src/db"
	"esnd/src/service"
	"esnd/src/util"
	"os"
	"strconv"
	"sync"
)

var Cfg *util.Config
var Service *service.NetService
var wg sync.WaitGroup

var configDefault = `[server]
service.port=3003
log.enable=true

[mysql]
db.user=esnd
db.addr=127.0.0.1:3306
db.pass=changeMe
db.database=esnd

[admin]
root.mask=changeMe`

func main() {
	util.SaySub("Main", "Protocol Version:"+strconv.Itoa(util.ProtocolVersion))
	configName := "config/esnd.conf"
	if len(os.Args) == 1 {

		fileExist, err := PathExists("config/esnd.conf")
		if err != nil {
			panic(err)
		}
		if !fileExist {
			util.SaySub("Main", "Config file not found,creating './config/esnd.conf' and exit.")

			exist, err := PathExists("config")
			if err != nil {
				panic(err)
			}
			if !exist {
				os.MkdirAll("config", os.ModePerm)
			}
			//create
			util.Sayln(configDefault)

			err = WriteFile("config/esnd.conf", configDefault)

			if err != nil {
				util.SaySub("Main", "Cannot create config file")
				panic(err)
			}

			util.SaySub("Main", "Config file created,please edit it,then restart program.")
			os.Exit(0)
		}
		configName = "config/esnd.conf"
	} else {
		configName = os.Args[1]
	}

	exist, err := PathExists(configName)
	if err != nil {
		panic(err)
	}
	if !exist {
		util.SaySub("Main", "Config file:"+configName+" not found.")
		os.Exit(1)
	}
	util.SaySub("Main", "Loading config file:"+configName)
	cfg, err := util.LoadConfig(configName)
	if err != nil {
		panic(err)
	}
	Cfg = cfg
	util.SaySub("Main", "Config file loaded.")

	if Cfg.GetAnyway("log.enable", "false") == "true" {
		util.EnableLog = true
		util.SaySub("Main", "See log in esnd.log.")
	}

	util.DebugMode, err = strconv.ParseBool(cfg.GetAnyway("debug.enable", "false"))
	if err != nil {
		util.SaySub("Main", "Config: debug.enable is not a bool")
	}
	util.SaySub("Main", "debugMode="+strconv.FormatBool(util.DebugMode))

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

func WriteFile(name string, str string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	write := bufio.NewWriter(file)
	write.WriteString(str)
	write.Flush()
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
