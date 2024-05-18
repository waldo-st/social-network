package port

import (
	"context"
	"social/internal/core/domain"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetCommentById(ctx context.Context, id int) (*domain.Comment, error)
	GetCommentsByUserId(ctx context.Context, id int) ([]*domain.Comment, error)
	GetCommentsByPostId(ctx context.Context, id int) ([]*domain.CommentInfo, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetCommentById(ctx context.Context, id int) (*domain.Comment, error)
	GetCommentsByUserId(ctx context.Context, id int) ([]*domain.Comment, error)
	GetCommentsByPostId(ctx context.Context, id int) ([]*domain.CommentInfo, error)
}
