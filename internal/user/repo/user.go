package repo

import (
	"context"
	"database/sql"
	"fmt"
	"image-gallery/internal/user/entity"
	"image-gallery/pkg/metrics"
	"log"
	"time"
)

func (r *Repository) CreateUser(user *entity.User) (*entity.User, error) {
	ok, fail := metrics.DatabaseQueryTime("Sign Up")
	defer fail()

	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int

	query := "INSERT INTO users(username, password, email, role, is_confirmed) VALUES ($1, $2, $3, 'client', false) returning id"
	err := r.db.QueryRowContext(c, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)

	if err != nil {
		return &entity.User{}, err
	}

	user.Id = int64(lastInsertId)
	ok()
	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	u := entity.User{}

	query := "SELECT id, email, username, password, role, is_confirmed FROM users WHERE email = $1"
	err := r.db.QueryRowContext(c, query, email).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role, &u.IsConfirmed)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *Repository) GetUserById(id int) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u := entity.User{}
	query := "SELECT id, email, username, password, role FROM users WHERE id = $1"

	err := r.db.QueryRowContext(c, query, id).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *Repository) GetUserByUsername(username string) (*entity.User, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u := entity.User{}
	query := "SELECT id, email, username, password, role FROM users WHERE username = $1"

	err := r.db.QueryRowContext(c, query, username).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role)
	if err != nil {
		return &entity.User{}, nil
	}

	return &u, nil
}

func (r *Repository) GetAllUsers() ([]*entity.User, error) {
	ok, fail := metrics.DatabaseQueryTime("Sign Up")
	defer fail()
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
			log.Fatalf("error in rows")
			return
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

	ok()
	return users, nil
}

func (r *Repository) DeleteUser(id int) error {
	ok, fail := metrics.DatabaseQueryTime("Sign Up")
	defer fail()

	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "delete from users where id=$1"
	queryAuth := "delete from user_token where user_id=$1"

	_, err := r.db.ExecContext(c, queryAuth, id)
	if err != nil {
		log.Fatalf("Can not delete the user Token: %s", err)
	}

	_, err = r.db.ExecContext(c, query, id)
	if err != nil {
		log.Fatalf("Can not delete the user: %s", err)
	}

	ok()
	return nil
}

func (r *Repository) UpdateUser(id int, newUsername string) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `update users set username= $1 where id = $2`
	_, err := r.db.ExecContext(c, query, newUsername, id)

	if err != nil {
		return fmt.Errorf("failed at query exec: %v", err)
	}

	return nil
}

func (r *Repository) UserCodeInsert(code *entity.UserCode) error {
	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var lastInsertId int

	query := `insert into "user_code" (user_id, user_code) values ($1, $2) returning id`
	err := r.db.QueryRowContext(c, query, code.UserId, code.UserCode).Scan(&lastInsertId)
	if err != nil {
		return fmt.Errorf("error in query row context: %s, UserId: %v", err, code.UserId)
	}

	code.Id = lastInsertId
	return nil
}

func (r *Repository) ConfirmUser(userCode string) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT user_id FROM user_code WHERE user_code = $1"
	var userId int
	err := r.db.QueryRowContext(c, query, userCode).Scan(&userId)
	if err != nil {
		log.Fatalf("Cannot retrieve user_id from user_code: %s", err)
		return err
	}

	updateQuery := "UPDATE users SET is_confirmed = true WHERE id = $1"
	_, err = r.db.ExecContext(c, updateQuery, userId)
	if err != nil {
		log.Fatalf("Cannot update users table: %s", err)
		return err
	}

	deleteQuery := "DELETE FROM user_code WHERE user_code = $1"
	_, err = r.db.ExecContext(c, deleteQuery, userCode)
	if err != nil {
		log.Fatalf("Cannot delete the user token from user_code: %s", err)
		return err
	}

	return nil
}
