package service

import (
	"context"
	"errors"
	"fmt"
	"social/internal/core/domain"
	"social/internal/core/port"
	"time"
)

type groupService struct {
	ur port.UserRepository
	// postRepo port.PostRepository
	// commentRepo port.CommentRepository
	gr      port.GroupRepository
	timeout time.Duration
}

func NewGroupService(ur port.UserRepository, gr port.GroupRepository) *groupService {
	return &groupService{ur, gr, time.Duration(3) * time.Second}
}

func (s *groupService) CreateGroup(ctx context.Context, gp *domain.Group) (*domain.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	user, err := s.ur.GetUserById(ctx, gp.Admin)
	if err != nil || user.Id == 0 {
		return nil, errors.New("invalid user")
	}

	g, err := s.gr.CreateGroup(ctx, gp)
	if err != nil {
		return nil, err
	}

	err = s.AddMember(ctx, g)
	return g, err
}

func (s *groupService) AddMember(ctx context.Context, g *domain.Group) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	fmt.Println("member top: ", g)
	if s.gr.IsMember(ctx, &domain.Member{UserId: g.Admin, GroupId: g.Id}) {
		return errors.New("add member service: already membber")
	}

	m := &domain.Member{UserId: g.Admin, GroupId: g.Id}

	err := s.gr.AddMember(ctx, m)
	fmt.Println("member bottom: ", g, err)
	return err
}

func (s *groupService) ListGroupMembers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.gr.ListGroupMembers(ctx, id)
}

func (s *groupService) ListGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.gr.ListGroups(ctx, id)
}

func (s *groupService) GetGroupsByUserId(ctx context.Context, id int) ([]*domain.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.gr.GetGroupsByUserId(ctx, id)
}

func (s *groupService) GetGroupById(ctx context.Context, id int) (*domain.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.gr.GetGroupById(ctx, id)
}

func (s *groupService) ListJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	return s.gr.GetJoinedGroups(ctx, id)
}

func (s *groupService) ListUnjoinedGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	return s.gr.GetUnjoinedGroups(ctx, id)
}
