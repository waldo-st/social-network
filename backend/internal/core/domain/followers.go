package domain

type Follow struct {
	FollowerId int `json:"followerid"`
	FolloweeId int `json:"followeeid"`
}
