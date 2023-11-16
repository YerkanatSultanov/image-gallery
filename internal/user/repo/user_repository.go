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

type Repository struct {
	db DBTX
}

type RepositoryInt interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserById(d int) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)
	DeleteUserByEmail() error
}

func NewRepository(db DBTX) *Repository {
	return &Repository{db: db}
}
