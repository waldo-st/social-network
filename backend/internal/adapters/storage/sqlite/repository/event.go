package repository

import (
	"context"
	"database/sql"
	"social/internal/core/domain"
	"sync"
)

type eventRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *eventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO event (GroupId, CreatorId, Title, Description, CreatedAt) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(ctx, event.GroupId, event.CreatorId, event.Title, event.Description, event.CreatedAt)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	event.Id = int(id)

	return event, nil
}

func (r *eventRepository) GetEventById(ctx context.Context, id int) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT id, groupId, creatorId, title, description, createdAt FROM event WHERE id=?`)
	if err != nil {
		return nil, err
	}

	var event = new(domain.Event)

	err = stm.QueryRowContext(ctx, id).Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.CreatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (r *eventRepository) GetEventsByCreatorId(ctx context.Context, id int) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT Id, GroupId, CreatorId, Title, Description, CreatedAt FROM Event WHERE CreatorId=?`)
	if err != nil {
		return nil, err
	}
	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var event = new(domain.Event)
	var eventsList []*domain.Event

	for rows.Next() {
		err := rows.Scan(event.Id, event.GroupId, event.CreatorId, event.Title, event.Description, event.CreatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}

func (r *eventRepository) ListGroupEvents(ctx context.Context, id int) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT Id, GroupId, CreatorId, Title, Description, CreatedAt FROM Event WHERE GroupId=?`)
	if err != nil {
		return nil, err
	}
	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var eventsList []*domain.Event

	for rows.Next() {
		var event = new(domain.Event)
		err := rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.CreatedAt)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}

// discuss
func (r *eventRepository) RespondToEvent(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO event_reaction (eventId, userId, status, issuedAt) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(ctx, reaction.EventId, reaction.UserId, reaction.Status, reaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	reaction.Id = int(id)

	return reaction, nil
}

// func (r *eventRepository) GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()

// 	stm, err := r.db.PrepareContext(ctx, `SELECT e.id, e.groupId, e.creatorId, e.title, e.description, e.createdAt FROM event e JOIN  event_reaction r ON e.id=r.eventId WHERE r.userId = ? AND r.status = ?`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rows, err := stm.QueryContext(ctx, id, "going")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var eventsList []*domain.Event

// 	for rows.Next() {
// 		var event = new(domain.Event)
// 		err := rows.Scan(&event.Id, &event.GroupId, &event.CreatorId, &event.Title, &event.Description, &event.CreatedAt)
// 		if err != nil && err != sql.ErrNoRows {
// 			return nil, err
// 		}
// 		eventsList = append(eventsList, event)
// 	}
// 	return eventsList, nil
// }
