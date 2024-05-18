package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"
	"time"
)

type EventService struct {
	userRepo  port.UserRepository
	groupRepo port.GroupRepository
	eventRepo port.EventRepository
	timeout   time.Duration
}

func NewEventService(ur port.UserRepository, gr port.GroupRepository, er port.EventRepository) *EventService {
	return &EventService{userRepo: ur, groupRepo: gr, eventRepo: er, timeout: time.Duration(3) * time.Second}
}

func (s *EventService) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if group, err := s.groupRepo.GetGroupById(ctx, event.GroupId); err != nil || group.Id == 0 {
		return nil, errors.New("group not found")
	}
	return s.eventRepo.CreateEvent(ctx, event)
}

func (s *EventService) GetEventById(ctx context.Context, id int) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	// make all needed validations here

	return nil, nil
}

func (s *EventService) GetEventsByCreatorId(ctx context.Context, id int) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	// make all needed validations here
	return nil, nil
}

func (s *EventService) ListGroupEvents(ctx context.Context, id int) ([]*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if group, err := s.groupRepo.GetGroupById(ctx, id); err != nil || group.Id == 0 {
		return nil, errors.New("group not found")
	}
	return s.eventRepo.ListGroupEvents(ctx, id)
}

func (s *EventService) RespondToEvent(ctx context.Context, reaction *domain.Reaction) (*domain.Reaction, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if event, err := s.eventRepo.GetEventById(ctx, reaction.EventId); err != nil || event.Id == 0 {
		return nil, errors.New("event not found")
	}

	return s.eventRepo.RespondToEvent(ctx, reaction)
}

// func (s *EventService) GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error) {
// 	return s.eventRepo.GetAttendingEvents(ctx, id)
// }
