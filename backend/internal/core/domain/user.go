package domain

import (
	"time"
)

type User struct {
	Id          int
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Avatar      string    `json:"avatar"`
	Nickname    string    `json:"nickname"`
	About       string    `json:"about"`
	IsPublic    bool      `json:"isPublic"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserLog struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRes struct {
	Id             int
	Email          string
	FirstName      string
	LastName       string
	DateOfBirth    time.Time
	Avatar         string
	Nickname       string
	About          string
	CreatedAt      time.Time
	IsPublic       bool
	IsFollowee     string
	NbrOfFollowers int
	NbrOfFollowee  int
	Token          string
}

type ProfilRes struct {
}
