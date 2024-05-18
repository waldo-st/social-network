package domain

import (
	"time"
)

type Notification struct {
	Id        int
	GroupId   int    `json:"gid"`
	Sender    int    `json:"sender"`
	Receiver  int    `json:"receiver"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt time.Time
}

type NotificationRes struct {
	Id       int    `json:"id"`
	GroupId  int    `json:"gid"`
	Username string `json:"username"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}
