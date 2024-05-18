package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"
	"time"
)

type followerService struct {
	r       port.FollowerRepository
	timeout time.Duration
}

func NewFollower(r port.FollowerRepository) *followerService {
	return &followerService{r, time.Duration(3) * time.Second}
}

func (fs *followerService) Follow(ctx context.Context, follow *domain.Follow) (*domain.Follow, error) {
	ctx, cancel := context.WithTimeout(ctx, fs.timeout)
	defer cancel()

	if fs.r.IsFollow(ctx, follow.FollowerId, follow.FolloweeId) {
		return nil, errors.New("you are already a follower")
	}

	f, err := fs.r.Add(ctx, follow)
	return f, err
}

func (fs *followerService) GetUserFollowers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, fs.timeout)
	defer cancel()

	return fs.r.ListUserFollowers(ctx, id)
}

func (fs *followerService) GetUserFollowee(ctx context.Context, id int) ([]*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, fs.timeout)
	defer cancel()

	return fs.r.ListUserFollowee(ctx, id)
}

func (fs *followerService) UnFollow(ctx context.Context, follower, followee int) error {
	ctx, cancel := context.WithTimeout(ctx, fs.timeout)
	defer cancel()
	err := fs.r.Remove(ctx, follower, followee)
	return err
}
