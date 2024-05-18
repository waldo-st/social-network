package port

import (
	"context"
	"social/internal/core/domain"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetPostById(ctx context.Context, id int) (*domain.Post, error)
	GetsPostsByUserId(ctx context.Context, id int) ([]*domain.Post, error)
	GetPostsByGroupId(ctx context.Context, id int) ([]*domain.Post, error)
	ListPosts(ctx context.Context, userId int) ([]*domain.PostInfo, error)
	SelectFollowers(ctx context.Context, followers []int, id int) error
}

type PostService interface {
	CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetPostById(ctx context.Context, id int) (*domain.Post, error)
	GetsPostsByUserId(ctx context.Context, id int) ([]*domain.Post, error)
	GetPostsByGroupId(ctx context.Context, id int) ([]*domain.Post, error)
	ListPosts(ctx context.Context, userId int) ([]*domain.PostInfo, error)
}
