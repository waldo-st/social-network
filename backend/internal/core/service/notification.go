package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"
)

type notificationService struct {
	repo port.NotificationRepository
}

func NewNotificationService(srv port.NotificationRepository) *notificationService {
	return &notificationService{srv}
}

func (n *notificationService) Add(ctx context.Context, notification *domain.Notification) (*domain.Notification, error) {
	_, ok := n.repo.Exist(ctx, notification)

	if !ok {
		return n.repo.Add(ctx, notification)
	}
	return nil, errors.New("notification: user has already request")
}

func (n *notificationService) Get(ctx context.Context, id int) (*domain.Notification, error) {
	return n.repo.Get(ctx, id)
}

func (n *notificationService) GetNotifications(ctx context.Context, receiver int) ([]*domain.NotificationRes, error) {
	return n.repo.GetNotifications(ctx, receiver)
}

func (n *notificationService) Exist(ctx context.Context, nt *domain.Notification) (*domain.NotificationRes, bool) {
	return n.repo.Exist(ctx, nt)
}

func (n *notificationService) Update(ctx context.Context, status string, notifid int) error {
	return n.repo.Update(ctx, status, notifid)
}

func (n *notificationService) Delete(ctx context.Context, sender, receiver int) error {
	return n.repo.Delete(ctx, sender, receiver)
}
