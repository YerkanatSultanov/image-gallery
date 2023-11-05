package repo

import (
	"context"
	"database/sql"
	"image-gallery/internal/user/entity"
)

func (r *repository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &entity.User{}, err
	}

	user.Id = int64(lastInsertId)
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	u := entity.User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *repository) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	u := entity.User{}
	query := "SELECT id, email, username, password FROM users WHERE id = $1"

	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.Id, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *repository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	query := "SELECT id, username, email FROM users"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle the error if needed.
		}
	}(rows)

	var users []*entity.User

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
