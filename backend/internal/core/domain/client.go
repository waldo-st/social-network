package domain //ws instead

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Chat
	Id       int    `json:"id"`
	RoomId   int    `json:"roomId"`
	Username string `json:"username"`
	Type     string `json:"type"`
}
