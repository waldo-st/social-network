package port

import (
	"context"
	"social/internal/core/domain"
)

type TokenRepository interface {
	Create(ctx context.Context, token *domain.Token) (string, error)
	Update(ctx context.Context, token *domain.Token) (string, error)
	Get(ctx context.Context, token string) (*domain.Token, error)
	GetTokenByUserId(ctx context.Context, id int) (*domain.Token, error)
	GetUserByToken(ctx context.Context, t string) (*domain.Token, error)
}

type TokenService interface {
	CreateToken(ctx context.Context, token *domain.User) (string, error)
	VerifyToken(ctx context.Context, token string) (int, error)
	Login(ctx context.Context, userReq *domain.UserLog) (*domain.User, error)
	GetTokenByUserId(ctx context.Context, id int) (*domain.Token, error)
	GetUserByToken(ctx context.Context, t string) (*domain.Token, error)
}
