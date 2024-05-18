package repository

import (
	"context"
	"database/sql"
	"html"
	"social/internal/core/domain"
	"sync"
)

type postRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO Post (userId, groupId,title, content, image, privacy, createdAt) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(
		ctx,
		post.UserId,
		post.GroupId,
		html.EscapeString(post.Title),
		html.EscapeString(post.Content),
		html.EscapeString(post.Image),
		post.Privacy,
		post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	post.Id = int(id)

	if err = r.SelectFollowers(ctx, post.SelectedUser, post.Id); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) GetPostById(ctx context.Context, id int) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stm, err := r.db.PrepareContext(ctx, `SELECT id, userId, title, content, image, privacy, createdAt FROM Post WHERE Id = ?`)
	if err != nil {
		return nil, err
	}

	var post = new(domain.Post)

	if err = stm.QueryRowContext(ctx, id).Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt); err != nil {
		return nil, err
	}

	return post, nil
}

func (r *postRepository) GetsPostsByUserId(ctx context.Context, id int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stm, err := r.db.PrepareContext(ctx, `SELECT id, userId, title, content, image, privacy, createdAt FROM Post WHERE userId = ?`)
	if err != nil {
		return nil, err
	}

	var postsList []*domain.Post

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post = new(domain.Post)
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt); err != nil && err != sql.ErrNoRows {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		postsList = append(postsList, post)
	}

	return postsList, nil
}

func (r *postRepository) GetPostsByGroupId(ctx context.Context, id int) ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT id, userId, groupId,title, content, image, privacy, createdAt FROM Post WHERE groupId = ?`)
	if err != nil {
		return nil, err
	}

	var postsList []*domain.Post

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post = new(domain.Post)
		if err := rows.Scan(&post.Id, &post.UserId, &post.GroupId, &post.Title, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt); err != nil && err != sql.ErrNoRows {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		postsList = append(postsList, post)
	}

	return postsList, nil
}

func (r *postRepository) ListPosts(ctx context.Context, id int) ([]*domain.PostInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`SELECT p.id, p.userId, p.title, p.content, p.image, p.privacy, p.createdAt, u.firstname, u.lastname, u.avatar,
			(SELECT count(*) FROM comment c  WHERE p.id = c.postId) as "commentsNbr"
		FROM Post p
		JOIN user u ON p.userId = u.id 
		WHERE p.privacy= 'public'
		UNION
		SELECT p.id, p.userId, p.title, p.content, p.image, p.privacy, p.createdAt, u.firstname, u.lastname, u.avatar,
			(SELECT count(*) FROM comment c  WHERE p.id = c.postId) as "commentsNbr"
		FROM Post p 
		JOIN user u ON p.userId = u.id
		JOIN follower f ON f.followeeId = p.userId 
		WHERE f.followerId = ? AND p.privacy= 'private'
		UNION
		SELECT p.id, p.userId, p.title, p.content, p.image, p.privacy, p.createdAt, u.firstname, u.lastname, u.avatar,
			(SELECT count(*) FROM comment c  WHERE p.id = c.postId) as "commentsNbr"
		FROM Post p 
		JOIN user u ON p.userId = u.id
		JOIN selected_user su ON su.postId = p.Id 
		WHERE p.privacy= 'almost private' AND su.followerId = ?
		ORDER BY p.id DESC
		`,
	)
	if err != nil {
		return nil, err
	}

	var postsList []*domain.PostInfo

	rows, err := stm.QueryContext(ctx, id, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post = new(domain.PostInfo)
		if err := rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt, &post.Username, &post.Lastname, &post.UserAvatar, &post.NbrComments); err != nil && err != sql.ErrNoRows {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		postsList = append(postsList, post)
	}

	return postsList, nil
}

func (r *postRepository) SelectFollowers(ctx context.Context, followers []int, id int) error {
	stm, err := r.db.PrepareContext(ctx, `INSERT INTO selected_user (postId, followerId) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	for _, f := range followers {
		_, err := stm.ExecContext(ctx, id, f)
		if err != nil {
			return err
		}

	}
	return nil
}
