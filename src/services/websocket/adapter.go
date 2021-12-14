package websocket

import (
	"encoding/json"
	"errors"
	"esnd/src/services"
	"esnd/src/util"
	"strconv"

	"github.com/gorilla/websocket"
)

type WSAdapter struct {
	Conn *websocket.Conn
}

type WSNetPackage struct {
	Code     int
	Crypto   int
	DataPack string
}

func (wsa *WSAdapter) HandShake() (bool, error) {
	//read netpackage
	np, err := wsa.Read()
	if err != nil {
		return false, err
	}
	//unmarshall json to PackTest interface{}
	pack := &services.PackTest{}
	err = json.Unmarshal([]byte(np.JSON), &pack)
	if err != nil {
		return false, err
	}

	//check
	if pack.Integer != 119812525 {
		return false, errors.New("invali identifier from ws connection:" + strconv.Itoa(pack.Integer))
	}
	pack.Integer = util.ProtocolVersion
	pack.Msg = "ProtocolVersion"
	_, err = wsa.Write(pack, 0)
	if err != nil {
		return false, errors.New("failed to reply protocol version:" + err.Error())
	}
	return true, nil
}

func (wsa *WSAdapter) Write(p interface{}, code int) (*services.NetPackage, error) {
	//convert data pkg to json
	jsonb, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	jsons := string(jsonb)
	//make wsnetpkg
	var wsnp WSNetPackage
	wsnp.Code = code
	wsnp.Crypto = 0
	wsnp.DataPack = jsons

	//send wsnetpkg
	err = wsa.Conn.WriteJSON(wsnp)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (wsa *WSAdapter) Read() (*services.NetPackage, error) {
	wsnp := &WSNetPackage{}
	err := wsa.Conn.ReadJSON(wsnp)
	if err != nil {
		return nil, err
	}
	var np services.NetPackage
	np.Code = wsnp.Code
	np.Crypto = wsnp.Crypto == 1
	np.Size = len([]byte(wsnp.DataPack))
	np.JSON = wsnp.DataPack
	if np.Crypto { //unsupported,return err when receive a encrypted package
		return nil, errors.New("encrypted package is now unsupported")
	}
	return &np, nil
}

func (wsa *WSAdapter) Dispose() {
	_ = wsa.Conn.Close()
}
