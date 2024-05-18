package repository

import (
	"context"
	"database/sql"
	"log"
	"social/internal/core/domain"
	"sync"
)

type notificationRepository struct {
	m  sync.RWMutex
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *notificationRepository {
	return &notificationRepository{db: db}
}

func (np *notificationRepository) Add(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	np.m.Lock()
	defer np.m.Unlock()

	stmt, err := np.db.PrepareContext(ctx, "INSERT INTO notification(sender, receiver, status, type, message,groupId) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	notif, err := stmt.ExecContext(ctx, n.Sender, n.Receiver, n.Status, n.Type, n.Message, n.GroupId)
	if err != nil {
		return nil, err
	}
	id, err := notif.LastInsertId()
	n.Id = int(id)
	return n, err
}
func (n *notificationRepository) Exist(ctx context.Context, nt *domain.Notification) (*domain.NotificationRes, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	stmt, err := n.db.PrepareContext(ctx, "SELECT id FROM notification WHERE sender=? AND receiver=? AND status != 'rejected' AND type=?;")
	if err != nil {
		return nil, false
	}

	var nr = new(domain.Notification)

	err = stmt.QueryRowContext(ctx, nt.Sender, nt.Receiver, nt.Type).Scan(&nr.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
	}

	return &domain.NotificationRes{Id: nr.Id}, true
}

func (n *notificationRepository) Get(ctx context.Context, id int) (*domain.Notification, error) {
	n.m.RLock()
	defer n.m.RUnlock()

	stm, err := n.db.PrepareContext(ctx, `SELECT id, sender, receiver, status, type, message, groupId FROM notification WHERE id = ?;`)
	if err != nil {
		return nil, err
	}

	var nt = new(domain.Notification)

	err = stm.QueryRowContext(ctx, id).Scan(&nt.Id, &nt.Sender, &nt.Receiver, &nt.Status, &nt.Type, &nt.Message, &nt.GroupId)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return nt, nil
}

func (n *notificationRepository) GetNotifications(ctx context.Context, receiver int) ([]*domain.NotificationRes, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	stmt, err := n.db.PrepareContext(ctx, "SELECT n.id, u.firstname, n.type, n.message, groupId FROM notification n JOIN user u ON n.sender = u.id WHERE receiver=? AND status='pending'")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, receiver)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	defer rows.Close()

	var notifications []*domain.NotificationRes
	for rows.Next() {
		var n = new(domain.NotificationRes)
		if err := rows.Scan(&n.Id, &n.Username, &n.Type, &n.Message, &n.GroupId); err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (n *notificationRepository) Update(ctx context.Context, status string, id int) error {
	n.m.Lock()
	defer n.m.Unlock()

	stmt, err := n.db.PrepareContext(ctx, "UPDATE notification SET status = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, status, id)
	return err
}

func (n *notificationRepository) Delete(ctx context.Context, sender, receiver int) error {
	n.m.Lock()
	defer n.m.Unlock()
	stm, err := n.db.PrepareContext(ctx, `DELETE FROM notification WHERE sender = ? AND receiver = ? AND status='accepted';`)
	if err != nil {
		return err
	}
	_, err = stm.ExecContext(ctx, sender, receiver)
	if err != nil {
		return err
	}

	return nil
}
