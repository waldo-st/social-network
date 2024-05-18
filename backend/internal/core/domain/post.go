package domain

import "time"

type Post struct {
	Id           int    `json:"id"`
	UserId       int    `json:"userid"`
	GroupId      int    `json:"groupid"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	Privacy      string `json:"privacy"`
	SelectedUser []int  `json:"selected_users"`
	CreatedAt    time.Time
}

type PostInfo struct {
	Id          int
	UserId      int
	Username    string
	Lastname    string
	UserAvatar  string
	Title       string
	Content     string
	Image       string
	Privacy     string
	CreatedAt   time.Time
	NbrComments int
}
