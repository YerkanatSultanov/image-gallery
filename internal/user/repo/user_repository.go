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
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}
