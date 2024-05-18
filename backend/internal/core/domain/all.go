package domain

type Response struct {
	Profile UserRes
	Users   []*UserRes
	Posts   []*PostInfo
	Events  []*Event
	Chats   []*Chat
	Token   string
}

type EventRes struct {
	Author string
	Info   Event
}
