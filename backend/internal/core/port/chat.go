package port

import (
	"context"
	"social/internal/core/domain"
)

type ChatRepository interface {
	CreateMessage(ctx context.Context, message *domain.Chat) (*domain.Chat, error)
	// GetChatsByIds(ctx context.Context, senderId int, receiverId int) ([]*domain.Chat, error)
	GetChatsByGroupId(ctx context.Context, id int) ([]*domain.Chat, error)
}

type ChatService interface {
	CreateMessage(ctx context.Context, message *domain.Chat) (*domain.Chat, error)
	// GetChatsByIds(ctx context.Context, senderId int, receiverId int) ([]*domain.Chat, error)
	GetChatsByGroupId(ctx context.Context, id int) ([]*domain.Chat, error)
	ReadMessage(ctx context.Context, client *domain.Client, hub *domain.Hub)
	WriteMessage(client *domain.Client)
}
