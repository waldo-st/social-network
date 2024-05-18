package port

import (
	"context"
	"social/internal/core/domain"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	ListGroupMembers(ctx context.Context, id int) ([]*domain.UserRes, error)
	ListGroups(ctx context.Context, id int) ([]*domain.Group, error)
	GetGroupById(ctx context.Context, id int) (*domain.Group, error)
	GetGroupsByUserId(ctx context.Context, id int) ([]*domain.Group, error)
	GetJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	GetUnjoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	AddMember(ctx context.Context, m *domain.Member) error
	IsMember(ctx context.Context, m *domain.Member) bool
	// invite and accept
}

type GroupService interface {
	CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error)
	ListGroupMembers(ctx context.Context, id int) ([]*domain.UserRes, error)
	ListGroups(ctx context.Context, id int) ([]*domain.Group, error)
	GetGroupById(ctx context.Context, id int) (*domain.Group, error)
	GetGroupsByUserId(ctx context.Context, id int) ([]*domain.Group, error)
	ListJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	ListUnjoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	AddMember(ctx context.Context, m *domain.Group) error
}
