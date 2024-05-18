package domain

import (
	"time"
)

type Group struct {
	Id          int
	Title       string
	Description string
	Admin       int
	Membership  string
	CreatedAt   time.Time
}
