package service

import (
	"context"
	"errors"
	"social/internal/core/domain"
	"social/internal/core/port"
	util "social/internal/core/utils"
	"time"

	"github.com/gofrs/uuid/v5"
)

type TokenService struct {
	userRepo     port.UserRepository
	postRepo     port.PostRepository
	commentRepo  port.CommentRepository
	groupRepo    port.GroupRepository
	followerRepo port.FollowerRepository
	// chatRepo     port.ChatRepository
	// notifRepo    port.NotificationRepository
	eventRepo port.EventRepository
	tokenRepo port.TokenRepository
	timeout   time.Duration
}

// func NewTokenService(ur port.UserRepository, pr port.PostRepository, cr port.CommentRepository, gr port.GroupRepository, fr port.FollowerRepository, chr port.ChatRepository, nr port.NotificationRepository, er port.EventRepository, tr port.TokenRepository) *TokenService {

func NewTokenService(ur port.UserRepository, pr port.PostRepository, cr port.CommentRepository, gr port.GroupRepository, fr port.FollowerRepository, er port.EventRepository, tr port.TokenRepository) *TokenService {
	return &TokenService{
		userRepo:     ur,
		postRepo:     pr,
		commentRepo:  cr,
		groupRepo:    gr,
		followerRepo: fr,
		// chatRepo: chr,
		// notifRepo: nr,
		eventRepo: er,
		tokenRepo: tr,
		timeout:   time.Duration(3) * time.Second,
	}
}

func (s *TokenService) Login(ctx context.Context, user *domain.UserLog) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.userRepo.GetUserByLogin(ctx, user)
	if err != nil || u.Id <= 0 {
		return nil, errors.New("invalid credential: username is not valid")
	}

	err = util.ComparePassword(user.Password, u.Password)
	if err != nil {
		return nil, errors.New("invalid credential: not a valid password")
	}

	return u, err
}

func (s *TokenService) CreateToken(ctx context.Context, user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	session, err := s.tokenRepo.GetTokenByUserId(ctx, user.Id)
	if err != nil {
		return "", errors.New("Unauthorized")
	}

	if session.UserId == user.Id {
		session.Ttl = time.Now().Add(time.Hour * 3 * 24)
		return s.tokenRepo.Update(ctx, session)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	token := &domain.Token{
		UserId: user.Id,
		Value:  id.String(),
		Ttl:    time.Now().Add(time.Hour * 3 * 24),
	}

	return s.tokenRepo.Create(ctx, token)
}

func (s *TokenService) GetUserByToken(ctx context.Context, t string) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.tokenRepo.GetUserByToken(ctx, t)
}

func (s *TokenService) GetTokenByUserId(ctx context.Context, id int) (*domain.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.tokenRepo.GetTokenByUserId(ctx, id)
}

func (s *TokenService) VerifyToken(ctx context.Context, value string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	token, err := s.tokenRepo.Get(ctx, value)
	if err != nil {
		return -1, err
	}
	if !token.Ttl.After(time.Now()) {
		return -1, errors.New("session has expired")
	}
	return token.UserId, nil
}

// func (s *TokenService) GetResponseData(ctx context.Context, user *domain.User) (*domain.Response, error) {
// 	ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()

// 	posts, err := s.postRepo.ListPosts(ctx, user.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	users, err := s.userRepo.ListUsers(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// events, err := s.eventRepo.GetAttendingEvents(ctx, user.Id)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	return &domain.Response{
// 		Users:  users,
// 		Posts:  posts,
// 		// Events: events,
// 	}, nil

// }
