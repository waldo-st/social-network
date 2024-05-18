package repository

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"social/internal/core/domain"
	"sync"
)

type groupRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *groupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(ctx context.Context, group *domain.Group) (*domain.Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO "group" (creatorId, title, description, createdAt) VALUES (?, ?, ?, ?);`)
	if err != nil {
		return nil, err
	}

	result, err := stm.ExecContext(ctx,
		group.Admin,
		html.EscapeString(group.Title),
		html.EscapeString(group.Description),
		group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	group.Id = int(id)

	return group, err
}

func (r *groupRepository) AddMember(ctx context.Context, m *domain.Member) error {
	// r.mu.Lock()
	// defer r.mu.Unlock()
	fmt.Println("add member repo")
	stm, err := r.db.PrepareContext(ctx, `INSERT INTO member (userId, groupId) VALUES (?, ?);`)
	if err != nil {
		return err
	}

	_, err = stm.ExecContext(ctx, m.UserId, m.GroupId)
	if err != nil {
		return err
	}

	return nil
}

func (r *groupRepository) IsMember(ctx context.Context, m *domain.Member) bool {
	// r.mu.RLock()
	// defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT userId, groupId FROM member WHERE userId=? AND groupId=?;`)
	if err != nil {
		return false
	}

	var M = new(domain.Member)

	err = stm.QueryRowContext(ctx, m.UserId, m.GroupId).Scan(&M.UserId, &M.GroupId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	fmt.Println("repo member: ", err)
	return true
}

func (r *groupRepository) ListGroupMembers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`SELECT id, email, firstName, lastName, dateOfBirth, avatar, nickname, about, isPublic FROM member m JOIN user u ON m.userId = u.id WHERE m.groupId=?;`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var usersList = []*domain.UserRes{}

	for rows.Next() {
		var user = new(domain.UserRes)
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic)

		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		usersList = append(usersList, user)
	}
	return nil, nil
}

func (r *groupRepository) ListGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	r.mu.RLock()

	stm, err := r.db.PrepareContext(
		ctx,
		`
		SELECT g.id, g.creatorId, g.title, g.description, g.createdAt,
		CASE 
			WHEN n.status = 'pending' THEN 'pending' 
			WHEN m.userId = ? AND m.groupId = g.id THEN 'true' 
			ELSE 'false'
		END AS membership
		FROM "group" g
		LEFT JOIN notification n ON n."type" = 'group'
		LEFT JOIN member m ON m.userId = ? AND m.groupId = g.id ORDER BY g.id
		`)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id, id)
	if err != nil {
		return nil, err
	}

	var groups []*domain.Group

	for rows.Next() {
		var group = new(domain.Group)
		if err := rows.Scan(&group.Id, &group.Admin, &group.Title, &group.Description, &group.CreatedAt, &group.Membership); err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		groups = append(groups, group)
	}
	r.mu.RUnlock()

	return groups, nil
}

func (r *groupRepository) GetGroupsByUserId(ctx context.Context, id int) ([]*domain.Group, error) {
	r.mu.RLock()

	stm, err := r.db.PrepareContext(
		ctx,
		`SELECT g.id, g.creatorId, g.title, g.description, g.createdAt FROM "group" g JOIN member m ON g.id=m.groupId WHERE m.userId=?`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var groups []*domain.Group

	for rows.Next() {
		var group = new(domain.Group)
		if err := rows.Scan(&group.Id, &group.Admin, &group.Title, &group.Description, &group.CreatedAt); err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
		groups = append(groups, group)
	}
	r.mu.RUnlock()

	return groups, nil
}

func (r *groupRepository) GetGroupById(ctx context.Context, id int) (*domain.Group, error) {
	r.mu.RLock()

	stm, err := r.db.PrepareContext(ctx, `SELECT id, creatorId, title, description, createdAt FROM "group" WHERE id=?`)
	if err != nil {
		return nil, err
	}

	var group = new(domain.Group)

	err = stm.QueryRowContext(ctx, id).Scan(&group.Id, &group.Admin, &group.Title, &group.Description, &group.CreatedAt)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r *groupRepository) GetJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`
		SELECT g.id, g.creatorId, g.title, g.description
        FROM "group" g
        JOIN member m ON g.id = m.groupId
        WHERE m.userId = ?
		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var groups []*domain.Group

	for rows.Next() {
		var group = new(domain.Group)
		err := rows.Scan(&group.Id, &group.Admin, &group.Title, &group.Description)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *groupRepository) GetUnjoinedGroups(ctx context.Context, id int) ([]*domain.Group, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`
		SELECT g.id, g.creatorId, g.title, g.description
		FROM "group" g
		WHERE g.id NOT IN (
			SELECT m.groupId
			FROM member m
			WHERE m.userId = ?
		);
		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	var groups []*domain.Group

	for rows.Next() {
		var group = new(domain.Group)
		err := rows.Scan(&group.Id, &group.Admin, &group.Title, &group.Description)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
