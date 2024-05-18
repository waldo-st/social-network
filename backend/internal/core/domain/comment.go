package domain

import "time"

type Comment struct {
	Id        int
	UserId    int
	PostId    int
	Content   string
	Image     string
	CreatedAt time.Time
}


type CommentInfo struct {
	Id          int
	UserId      int
	Username    string
	Lastname    string
	UserAvatar  string
	Title       string
	Content     string
	Image       string
	CreatedAt   time.Time
}
