package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if websocket.IsWebSocketUpgrade(r) {
		wsc, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			panic(err)
		}
		_, bytes, err := wsc.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	} else {
		return
	}
}
func main() {
	http.HandleFunc("/ws", ServeHTTP)
	http.ListenAndServe(":3003", nil)
}
