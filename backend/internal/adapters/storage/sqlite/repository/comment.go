package repository

import (
	"context"
	"database/sql"
	"social/internal/core/domain"
	"sync"
)

type commentRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO Comment (UserId, PostId, Content, Image, CreatedAt) VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(ctx, comment.UserId, comment.PostId, comment.Content, comment.Image, comment.CreatedAt)

	if err != nil {
		return nil, err
	}

	if id, err := result.LastInsertId(); err == nil {
		comment.Id = int(id)
		return comment, err
	}

	return comment, nil
}

func (r *commentRepository) GetCommentById(ctx context.Context, id int) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stm, err := r.db.PrepareContext(ctx, `SELECT Id, UserId, PostId, Content, Image, CreatedAt FROM Comment WHERE Id = ?`)
	if err != nil {
		return nil, err
	}

	var comment = new(domain.Comment)

	if err = stm.QueryRowContext(ctx, id).Scan(&comment.Id, &comment.UserId, &comment.Content, &comment.Image, &comment.CreatedAt); err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *commentRepository) GetCommentsByUserId(ctx context.Context, id int) ([]*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT Id, UserId, PostId, Content, Image, CreatedAt FROM Comment WHERE UserId = ?`)
	if err != nil {
		return nil, err
	}

	var comment = new(domain.Comment)

	var commentList []*domain.Comment

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&comment.Id, &comment.UserId, &comment.Content, &comment.Image, &comment.CreatedAt); err != nil && err != sql.ErrNoRows {
			return nil, err
		}
	}

	return commentList, nil
}

func (r *commentRepository) GetCommentsByPostId(ctx context.Context, id int) ([]*domain.CommentInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`SELECT c.id, c.userId, c.content, c.image, c.createdAt, u.firstname, u.lastname, u.avatar
		FROM comment c
		JOIN user u ON c.userId = u.id
		WHERE c.postId = ?`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var commentList []*domain.CommentInfo
	for rows.Next() {
		var comment = new(domain.CommentInfo)
		if err := rows.Scan(&comment.Id, &comment.UserId, &comment.Content, &comment.Image, &comment.CreatedAt, &comment.Username, &comment.Lastname, &comment.UserAvatar); err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		commentList = append(commentList, comment)
	}

	return commentList, nil
}
