package domain

import "time"

type Token struct {
	Id       int
	UserId   int
	Value    string
	Username string
	Ttl      time.Time
}
