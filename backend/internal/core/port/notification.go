package port

import (
	"context"
	"social/internal/core/domain"
)

type NotificationRepository interface {
	Get(ctx context.Context, id int) (*domain.Notification, error)
	Update(ctx context.Context, status string, notifid int) error
	Add(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
	Exist(ctx context.Context, nt *domain.Notification) (*domain.NotificationRes, bool)
	GetNotifications(ctx context.Context, receiver int) ([]*domain.NotificationRes, error)
	Delete(ctx context.Context, sender, receiver int) error
}

type NotificationService interface {
	Get(ctx context.Context, id int) (*domain.Notification, error)
	Update(ctx context.Context, status string, notifid int) error
	Add(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
	Exist(ctx context.Context, nt *domain.Notification) (*domain.NotificationRes, bool)
	GetNotifications(ctx context.Context, receiver int) ([]*domain.NotificationRes, error)
	Delete(ctx context.Context, sender, receiver int) error
}
