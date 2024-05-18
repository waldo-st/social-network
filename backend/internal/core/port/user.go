package port

import (
	"context"
	"social/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserById(ctx context.Context, id int) (*domain.UserRes, error)
	GetUserByLogin(ctx context.Context, user *domain.UserLog) (*domain.User, error)
	GetOwnPosts(ctx context.Context, id int) ([]*domain.PostInfo, error)
	ListUsers(ctx context.Context, id int) ([]*domain.UserRes, error)
	GetConnections(ctx context.Context, id int) ([]*domain.UserRes, error)
	GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error)
	// GetJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	IsFollowee(ctx context.Context, userId int, followeeId int) (*domain.Follow, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type UserService interface {
	Register(ctx context.Context, user *domain.User) error
	Login(ctx context.Context, userReq *domain.UserLog) (*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.UserRes, error)
	GetOwnPosts(ctx context.Context, id int) ([]*domain.PostInfo, error)
	ListUsers(ctx context.Context, id int) ([]*domain.UserRes, error)
	ListConnections(ctx context.Context, id int) ([]*domain.UserRes, error)
	ListAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error)
	// ListJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error)
	CanSeeProfil(ctx context.Context, userId int, connection *domain.UserRes) (bool, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.UserRes, error)
	DeleteUser(ctx context.Context, id int) error
}
