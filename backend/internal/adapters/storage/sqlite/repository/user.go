package repository

import (
	"context"
	"database/sql"
	"html"
	"social/internal/core/domain"
	"sync"
)

type userRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	r.mu.Lock()

	stm, err := r.db.PrepareContext(ctx, `INSERT INTO user (email, password, firstName, lastName, dateOfBirth, avatar, nickname, about, createdAt, isPublic) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = stm.ExecContext(ctx,
		user.Email,
		user.Password,
		html.EscapeString(user.FirstName),
		html.EscapeString(user.LastName),
		user.DateOfBirth,
		html.EscapeString(user.Avatar),
		html.EscapeString(user.Nickname),
		html.EscapeString(user.About),
		user.CreatedAt,
		user.IsPublic,
	)
	if err != nil {
		return err
	}

	r.mu.Unlock()
	return nil
}

func (r *userRepository) ListUsers(ctx context.Context, id int) ([]*domain.UserRes, error) {
	r.mu.RLock()

	stm, err := r.db.PrepareContext(
		ctx,
		`
		SELECT u.id, u.email, u.firstname, u.lastname, u.dateofBirth, u.avatar, u.nickname, u.about, u.isPublic,
		CASE 
			WHEN n.status = 'pending' THEN 'pending' 
			WHEN f."followerId" = ? AND f."followeeId" = u.id THEN 'true' 
			ELSE 'false'
		END AS connection
		FROM user u
		LEFT JOIN notification n ON u.id = n.receiver AND n.sender = ? AND n.type = 'follow' 
		LEFT JOIN follower f ON f.followeeId = u.id AND f.followerId = ?
		WHERE u.id != ? ORDER BY u.id

		`)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id, id, id, id)
	if err != nil {
		return nil, err
	}

	var usersList = []*domain.UserRes{}

	for rows.Next() {
		var user = new(domain.UserRes)
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic, &user.IsFollowee)

		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		usersList = append(usersList, user)
	}
	r.mu.RUnlock()

	return usersList, nil
}

func (r *userRepository) GetUserByLogin(ctx context.Context, u *domain.UserLog) (*domain.User, error) {
	r.mu.RLock()
	stm, err := r.db.PrepareContext(ctx, `SELECT id, email, password, firstName, lastName, dateOfBirth, avatar, nickname, about,isPublic FROM user WHERE email = ?;`)
	if err != nil {
		return nil, err
	}

	var user = new(domain.User)

	err = stm.QueryRowContext(ctx, u.Login).Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	r.mu.RUnlock()
	return user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*domain.UserRes, error) {
	r.mu.RLock()
	stm, err := r.db.PrepareContext(ctx, `SELECT id, email, firstName, lastName, dateOfBirth, avatar, nickname, about, createdAt, isPublic FROM user WHERE id = ?;`)
	if err != nil {
		return nil, err
	}

	var user = new(domain.UserRes)

	err = stm.QueryRowContext(ctx, id).Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.CreatedAt, &user.IsPublic)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	r.mu.RUnlock()
	return user, nil
}

func (r *userRepository) GetOwnPosts(ctx context.Context, id int) ([]*domain.PostInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stm, err := r.db.PrepareContext(
		ctx,
		`SELECT p.id, p.userId, p.title, p.content, p.image, p.privacy, p.createdAt, u.firstname, u.lastname, u.avatar,
			(SELECT count(*) FROM comment c  WHERE p.id = c.postId) as "commentsNbr"
		FROM Post p 
		JOIN user u ON p.userId = u.id 
		WHERE p.userId = ?
		`)
	if err != nil {
		return nil, err
	}

	var postsList []*domain.PostInfo

	rows, err := stm.QueryContext(ctx, id)
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

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	stm, err := r.db.PrepareContext(ctx, `UPDATE user SET email = ?, password = ?, firstName = ?, lastName = ?, dateOfBirth = ?, avatar = ?, nickname = ?, about = ?, createdAt = ?, isPublic = ? WHERE id = ?;`)
	if err != nil {
		return nil, err
	}
	_, err = stm.ExecContext(ctx, user.Email, user.Password, user.FirstName, user.LastName, user.DateOfBirth, user.Avatar, user.Nickname, user.About, user.CreatedAt, user.IsPublic, user.Id)
	if err != nil {
		return nil, err
	}
	r.mu.Unlock()
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	r.mu.Lock()
	stm, err := r.db.PrepareContext(ctx, `DELETE FROM user WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stm.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	r.mu.Unlock()
	return nil
}

func (r *userRepository) GetAttendingEvents(ctx context.Context, id int) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT e.id, e.groupId, e.creatorId, e.title, e.description, e.createdAt FROM event e JOIN  event_reaction r ON e.id=r.eventId WHERE r.userId = ? AND r.status = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stm.QueryContext(ctx, id, "going")
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

func (r *userRepository) GetConnections(ctx context.Context, id int) ([]*domain.UserRes, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(
		ctx,
		`
		SELECT DISTINCT u.id, u.email, u.firstName, u.lastName, u.dateOfBirth, u.avatar, u.nickname, u.about, u.isPublic
        FROM user u JOIN follower f ON u.id = f.followerId
        WHERE f.followeeId = ?
        
        UNION
        
        SELECT DISTINCT u.id, u.email, u.firstName, u.lastName, u.dateOfBirth, u.avatar, u.nickname, u.about, u.isPublic
        FROM user u JOIN follower f ON u.id = f.followeeId
        WHERE f.followerId = ?

		`,
	)
	if err != nil {
		return nil, err
	}

	rows, err := stm.QueryContext(ctx, id, id)
	if err != nil {
		return nil, err
	}

	var connections []*domain.UserRes

	for rows.Next() {
		var user = new(domain.UserRes)
		err := rows.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.About, &user.IsPublic)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		connections = append(connections, user)
	}

	return connections, nil
}

func (r *userRepository) IsFollowee(ctx context.Context, userId int, followeeId int) (*domain.Follow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stm, err := r.db.PrepareContext(ctx, `SELECT followerId, followeeId FROM follower WHERE followerId=? AND followeeId=?`)
	if err != nil {
		return nil, err
	}
	var follower = new(domain.Follow)

	err = stm.QueryRowContext(ctx, userId, followeeId).Scan(&follower.FollowerId, &follower.FolloweeId)
	if err != nil {
		return nil, err
	}
	return follower, nil
}

// // func (r *userRepository) GetJoinedGroups(ctx context.Context, id int) ([]*domain.Group, error) {
// // 	r.mu.RLock()
// // 	defer r.mu.RUnlock()

// // 	stm, err := r.db.PrepareContext(
// // 		ctx,
// // 		`
// // 		SELECT g.id, g.creatorId, g.title, g.description
// //         FROM "group" g
// //         JOIN member m ON g.id = m.groupId
// //         WHERE m.userId = ?
// // 		`,
// // 	)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	rows, err := stm.QueryContext(ctx, id)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	var groups []*domain.Group

// // 	for rows.Next() {
// // 		var group = new(domain.Group)
// // 		err := rows.Scan(&group.Id, &group.Admin, &group.Title, &group.Description)
// // 		if err != nil && err != sql.ErrNoRows {
// // 			return nil, err
// // 		}
// // 		groups = append(groups, group)
// // 	}

// // 	return groups, nil
// // }
