package repository

import (
	"context"
	"database/sql"
	"html"
	"social/internal/core/domain"
	"sync"
)

type chatRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *chatRepository {
	return &chatRepository{db: db}
}

// GetMessagesByIds(ctx context.Context, senderId int, receiverId int) ([]*domain.Chat, error)

func (r *chatRepository) CreateMessage(ctx context.Context, message *domain.Chat) (*domain.Chat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO chat (senderId, username, groupId, content, image, createdAt) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(
		ctx,
		message.SenderId,
		html.EscapeString(message.Username),
		message.GroupId,
		html.EscapeString(message.Content),
		html.EscapeString(message.Image),
		message.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	message.Id = int(id)

	return message, nil
}

func (r *chatRepository) GetChatsByGroupId(ctx context.Context, id int) ([]*domain.Chat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT id, senderId, username, groupId, content, image, createdAt FROM chat WHERE chat.groupId = ?`)
	if err != nil {
		return nil, err
	}

	var messages []*domain.Chat

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg = new(domain.Chat)
		err = rows.Scan(&msg.Id, &msg.SenderId, &msg.Username, &msg.GroupId, &msg.Content, &msg.Image, &msg.CreatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
