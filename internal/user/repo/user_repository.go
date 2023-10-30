package repo

import (
	"context"
	"database/sql"
	"image-gallery/internal/user/entity"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

type Repository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserById(ctx context.Context, id int) (*entity.User, error)
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

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
