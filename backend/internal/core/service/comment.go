package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"

	"time"
)

type commentService struct {
	userRepo    port.UserRepository
	postRepo    port.PostRepository
	commentRepo port.CommentRepository
	timeout     time.Duration
}

func NewCommentService(ur port.UserRepository, pr port.PostRepository, cr port.CommentRepository) *commentService {
	return &commentService{userRepo: ur, postRepo: pr, commentRepo: cr, timeout: time.Duration(3) * time.Second}
}

func (s *commentService) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.userRepo.GetUserById(ctx, comment.UserId)
	if err != nil || user.Id == 0 {
		return nil, errors.New("No user associated with this comment")
	}

	post, err := s.postRepo.GetPostById(ctx, comment.PostId)
	if err != nil || post.Id == 0 {
		return nil, errors.New("no post associated with this comment")
	}

	comment.CreatedAt = time.Now()
	comment, err = s.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetCommentById(ctx context.Context, id int) (*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	comment, err := s.commentRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetCommentsByUserId(ctx context.Context, id int) ([]*domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	comment, err := s.commentRepo.GetCommentsByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) GetCommentsByPostId(ctx context.Context, id int) ([]*domain.CommentInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	comment, err := s.commentRepo.GetCommentsByPostId(ctx, id)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
