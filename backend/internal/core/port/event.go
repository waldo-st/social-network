package port

import (
	"context"
	"social/internal/core/domain"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetEventById(ctx context.Context, id int) (*domain.Event, error)
	// GetEventsByCreatorId(ctx context.Context, id int) ([]*domain.Event, error)
	ListGroupEvents(ctx context.Context, id int) ([]*domain.Event, error)
	RespondToEvent(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	// GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error)
}

type EventService interface {
	CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetEventById(ctx context.Context, id int) (*domain.Event, error)
	// GetEventsByCreatorId(ctx context.Context, id int) ([]*domain.Event, error)
	ListGroupEvents(ctx context.Context, id int) ([]*domain.Event, error)
	RespondToEvent(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error)
	// GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error)
}
