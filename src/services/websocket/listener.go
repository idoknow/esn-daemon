package websocket

import (
	"esnd/src/services"
	"esnd/src/util"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type WSService struct {
	Port int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func MakeService(port int) (*WSService, error) {
	var wss WSService
	wss.Port = port
	http.HandleFunc("/ws", AcceptHTTP)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
	return &wss, nil
}

/*
Be called each time when a connection incoming
*/
func AcceptHTTP(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		wsc, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			util.SaySub("Listener", "err:While accepting websocket connection:"+err.Error())
			return
		}
		go makeWSHandler(wsc)
	}
}

func makeWSHandler(wsc *websocket.Conn) {
	wsa := &WSAdapter{
		wsc,
	}
	hsResult, err := wsa.HandShake()
	if err != nil {
		util.SaySub("Listener", "err:While ws handshaking:"+err.Error())
		return
	}
	if !hsResult {
		util.SaySub("Listener", "Failed to ws handshake")
		return
	}
	services.CreateHandler(wsa)
}
