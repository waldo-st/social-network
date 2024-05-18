package service

import (
	"context"
	"errors"

	"social/internal/core/domain"
	"social/internal/core/port"
	util "social/internal/core/utils"
	"time"
)

type userService struct {
	repo    port.UserRepository
	timeout time.Duration
}

func NewUserService(r port.UserRepository) *userService {
	return &userService{r, time.Duration(3) * time.Second}
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := util.CheckFields(user); err != nil {
		return err
	}

	pwd, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pwd
	err = s.repo.CreateUser(ctx, user)
	return err
}

func (s *userService) Login(ctx context.Context, user *domain.UserLog) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	u, err := s.repo.GetUserByLogin(ctx, user)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	err = util.ComparePassword(user.Password, u.Password)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return u, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.repo.GetUserById(ctx, id)
	return user, err
}

func (s *userService) GetOwnPosts(ctx context.Context, id int) ([]*domain.PostInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	posts, err := s.repo.GetOwnPosts(ctx, id)
	return posts, err
}

func (s *userService) ListUsers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.repo.ListUsers(ctx, id)
	return user, err
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) (*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.repo.UpdateUser(ctx, user)
	response := &domain.UserRes{
		Id:          user.Id,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Avatar:      user.Avatar,
		Nickname:    user.Nickname,
		About:       user.About,
		CreatedAt:   user.CreatedAt,
		IsPublic:    user.IsPublic,
	}
	return response, err
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	err := s.repo.DeleteUser(ctx, id)
	return err
}

func (s *userService) ListAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error) {
	return s.repo.GetAttendingEvents(ctx, id)
}

func (s *userService) ListConnections(ctx context.Context, id int) ([]*domain.UserRes, error) {
	return s.repo.GetConnections(ctx, id)
}

func (s *userService) CanSeeProfil(ctx context.Context, userId int, connection *domain.UserRes) (bool, error) {
	if connection.IsPublic {
		return true, nil
	}

	if userId == connection.Id {
		return true, nil
	}

	follower, err := s.repo.IsFollowee(ctx, userId, connection.Id)
	if err != nil {
		return false, nil
	}
	if follower.FolloweeId != connection.Id {
		return false, nil
	}
	return true, nil
}
