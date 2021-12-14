package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

var url = "ws://localhost:3004/ws"

type WSNetPackage struct {
	Code     int
	Crypto   int
	DataPack string
}

type PackTest struct { //0 both
	Integer int
	Msg     string
	Token   string
}
type PackLogin struct { //1 client
	User  string
	Pass  string
	Token string
}
type PackPush struct { //3 client
	Target   string
	Time     string
	Title    string
	Content  string
	Token    string
	Realtime bool
}

var handshakeDataPkg = &PackTest{
	119812525,
	"",
	"handshake",
}

var loginDataPkg = &PackLogin{
	"root",
	"changeMe",
	"login",
}

var pushDataPkg = &PackPush{
	"_global_",
	"20211214-10:58",
	"FirstNotificationPushedThrWebSocket",
	"RockChin is a genius.",
	"pushNotification",
	false,
}

var wsnp = &WSNetPackage{
	0,
	0,
	"",
}

func main() {

	//handshake
	jsonb, err := json.Marshal(handshakeDataPkg)
	if err != nil {
		panic(err)
	}
	wsnp.DataPack = string(jsonb)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}
	err = conn.WriteJSON(wsnp)
	if err != nil {
		panic(err)
	}

	protocol := &WSNetPackage{}

	err = conn.ReadJSON(protocol)
	if err != nil {
		panic(err)
	}

	fmt.Println("ProtocolVersion Reply:" + protocol.DataPack)

	//send PackLogin

	wsnp.Code = 1
	loginJSON, err := json.Marshal(loginDataPkg)
	if err != nil {
		panic(err)
	}
	wsnp.DataPack = string(loginJSON)

	err = conn.WriteJSON(wsnp)
	if err != nil {
		panic(err)
	}

	handleResult(conn)

	//push notification
	wsnp.Code = 3
	pushJSON, err := json.Marshal(pushDataPkg)
	if err != nil {
		panic(err)
	}
	wsnp.DataPack = string(pushJSON)

	err = conn.WriteJSON(wsnp)
	if err != nil {
		panic(err)
	}

	handleResult(conn)
}

func handleResult(conn *websocket.Conn) {

	result := &WSNetPackage{}

	err := conn.ReadJSON(result)
	if err != nil {
		panic(err)
	}
	fmt.Println("LoginResult:" + result.DataPack)
}
