package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		if err := conn.WriteMessage(msgType, msg); err != nil {
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("listening on 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
