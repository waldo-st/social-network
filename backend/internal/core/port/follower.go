package port

import (
	"context"
	"social/internal/core/domain"
)

type FollowerRepository interface {
	Add(ctx context.Context, f *domain.Follow) (*domain.Follow, error)
	IsFollow(ctx context.Context, sender, receiver int) bool
	ListUserFollowers(ctx context.Context, id int) ([]*domain.UserRes, error)
	ListUserFollowee(ctx context.Context, id int) ([]*domain.UserRes, error)
	Remove(ctx context.Context, follower, followee int) error
}

type FollowerService interface {
	Follow(ctx context.Context, f *domain.Follow) (*domain.Follow, error)
	GetUserFollowers(ctx context.Context, id int) ([]*domain.UserRes, error)
	GetUserFollowee(ctx context.Context, id int) ([]*domain.UserRes, error)
	UnFollow(ctx context.Context, follower, followee int) error
}
