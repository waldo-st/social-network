package domain

import "time"

type Chat struct {
	Id        int
	SenderId  int
	Username  string
	GroupId   int
	Content   string
	Image     string
	Type      string
	CreatedAt time.Time
}

type WsChatReq struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
