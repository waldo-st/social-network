package repository

import (
	"context"
	"database/sql"
	"fmt"
	"social/internal/core/domain"
	"sync"
)

type tokenRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *tokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *domain.Token) (string, error) {
	fmt.Println("Token:", token)
	stm, err := r.db.PrepareContext(ctx, `INSERT INTO session (userId, cookie ,ttl) VALUES(?, ?, ?)`)
	if err != nil {
		return "", err
	}
	_, err = stm.ExecContext(ctx, token.UserId, token.Value, token.Ttl)
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

func (r *tokenRepository) Update(ctx context.Context, token *domain.Token) (string, error) {
	stm, err := r.db.PrepareContext(ctx, `UPDATE session SET userId=?, cookie=? ,Ttl=? WHERE userId=?`)
	if err != nil {
		return "", err
	}
	_, err = stm.ExecContext(ctx, token.UserId, token.Value, token.Ttl, token.UserId)
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

func (r *tokenRepository) GetTokenByUserId(ctx context.Context, id int) (*domain.Token, error) {
	token := new(domain.Token)
	stm, err := r.db.PrepareContext(ctx, `SELECT userId, cookie, ttl FROM session WHERE userId=?`)
	if err != nil {
		return nil, err
	}
	err = stm.QueryRowContext(ctx, id).Scan(&token.UserId, &token.Value, &token.Ttl)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return token, nil
}

func (r *tokenRepository) GetUserByToken(ctx context.Context, t string) (*domain.Token, error) {
	token := new(domain.Token)
	stm, err := r.db.PrepareContext(ctx, `SELECT s.userId, s.ttl, u.firstname FROM session s JOIN user u ON s.userId = u.id WHERE cookie=?`)
	if err != nil {
		return nil, err
	}

	err = stm.QueryRowContext(ctx, t).Scan(&token.UserId, &token.Ttl, &token.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return token, nil
}

func (r *tokenRepository) Get(ctx context.Context, value string) (*domain.Token, error) {
	token := new(domain.Token)
	stm, err := r.db.PrepareContext(ctx, `SELECT userId, ttl FROM session WHERE cookie=?`)
	if err != nil {
		return nil, err
	}
	err = stm.QueryRowContext(ctx, value).Scan(&token.UserId, &token.Ttl)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return token, nil
}
