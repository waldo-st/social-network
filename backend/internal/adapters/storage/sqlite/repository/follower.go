package repository

import (
	"context"
	"database/sql"
	"social/internal/core/domain"
	"sync"
)

type followerRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *followerRepository {
	return &followerRepository{db: db}
}

func (r *followerRepository) Add(ctx context.Context, f *domain.Follow) (*domain.Follow, error) {
	r.mu.Lock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO follower (followerid, followeeid) VALUES (?, ?);`)
	if err != nil {
		return nil, err
	}

	_, err = stm.ExecContext(ctx, f.FollowerId, f.FolloweeId)
	if err != nil {
		return nil, err
	}

	r.mu.Unlock()
	return f, nil
}

func (n *followerRepository) IsFollow(ctx context.Context, sender, receiver int) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	stm, err := n.db.PrepareContext(ctx, `SELECT id FROM follower WHERE followerId=? AND followeeId=?;`)
	if err != nil {
		return false
	}

	var id int
	err = stm.QueryRowContext(ctx, sender, receiver).Scan(&id)
	return err == nil
}

func (r *followerRepository) ListUserFollowers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT u.id, u.email, u.firstName, u.lastName, u.dateOfBirth, u.avatar, u.nickname, u.about, u.isPublic FROM user u JOIN follower f ON u.id = f.followerId WHERE f.followeeId = ?;`)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var followers []*domain.UserRes

	for rows.Next() {
		var user = new(domain.UserRes)
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic)

		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		followers = append(followers, user)
	}

	return followers, nil
}

func (r *followerRepository) ListUserFollowee(ctx context.Context, id int) ([]*domain.UserRes, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT u.id, u.email, u.firstName, u.lastName, u.dateOfBirth, u.avatar, u.nickname, u.about, u.isPublic FROM user u JOIN follower f ON u.id = f.followeeId WHERE f.followerId = ?;`)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var followee []*domain.UserRes

	for rows.Next() {
		var user = new(domain.UserRes)
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic)

		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		followee = append(followee, user)
	}

	return followee, nil
}

func (r *followerRepository) Remove(ctx context.Context, follower, followee int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	stm, err := r.db.PrepareContext(ctx, `DELETE FROM follower WHERE followerid = ? AND followeeid = ?;`)
	if err != nil {
		return err
	}
	_, err = stm.ExecContext(ctx, follower, followee)
	if err != nil {
		return err
	}

	return nil
}
