package repo

import (
	"context"
	"database/sql"
	"image-gallery/internal/user/entity"
	"log"
	"time"
)

func (r *Repository) CreateUser(user *entity.User) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(c, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		return &entity.User{}, err
	}

	user.Id = int64(lastInsertId)
	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	u := entity.User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(c, query, email).Scan(&u.Id, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *Repository) GetUserById(id int) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u := entity.User{}
	query := "SELECT id, email, username, password FROM users WHERE id = $1"

	err := r.db.QueryRowContext(c, query, id).Scan(&u.Id, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *Repository) GetAllUsers() ([]*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, username, email FROM users"
	rows, err := r.db.QueryContext(c, query)
	if err != nil {
		log.Fatalf("Error in database: %s", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

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

func (r *Repository) DeleteUserByEmail(email string) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "delete from users where email=$1"

	err := r.db.QueryRowContext(c, query, email)
	if err != nil {
		log.Fatalf("Can not delete the user: %s", err)
	}

	return nil
}