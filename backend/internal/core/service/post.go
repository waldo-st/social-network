package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"
	"strings"
	"time"
)

type postService struct {
	userRepo port.UserRepository
	postRepo port.PostRepository
	timeout  time.Duration
}

func NewPostService(ur port.UserRepository, pr port.PostRepository) *postService {
	return &postService{userRepo: ur, postRepo: pr, timeout: time.Duration(3) * time.Second}
}

func (s *postService) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if strings.TrimSpace(post.Title) == "" || strings.TrimSpace(post.Content) == "" {
		return nil, errors.New("post title or content cannot be empty")
	}

	user, err := s.userRepo.GetUserById(ctx, post.UserId)
	if err != nil || user.Id == 0 {
		return nil, errors.New("No user associated with this post")
	}

	if post.GroupId != 0 {
		post.Privacy = "public"
	}

	post, err = s.postRepo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) GetPostById(ctx context.Context, id int) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	post, err := s.postRepo.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) GetsPostsByUserId(ctx context.Context, id int) ([]*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	posts, err := s.postRepo.GetsPostsByUserId(ctx, id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *postService) GetPostsByGroupId(ctx context.Context, id int) ([]*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	posts, err := s.postRepo.GetPostsByGroupId(ctx, id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *postService) ListPosts(ctx context.Context, id int) ([]*domain.PostInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	posts, err := s.postRepo.ListPosts(ctx, id)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
