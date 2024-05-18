package util

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func Broadcast(ws *websocket.Conn, notif interface{}) error {
	msg, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		return err
	}
	fmt.Println("Message Broadcast => ", string(msg))
	return nil
}
