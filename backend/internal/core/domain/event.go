package domain

import "time"

type Event struct {
	Id          int
	GroupId     int
	CreatorId   int
	Title       string
	Description string
	Option      string
	CreatedAt   time.Time
}

type Reaction struct {
	Id        int
	EventId   int
	UserId    int
	Status    string
	CreatedAt time.Time
}

type EventReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type EventReaction struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}
