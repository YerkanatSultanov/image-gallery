package repo

import (
	"context"
	"database/sql"
	"image-gallery/internal/auth/entity"
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
	CreateUserToken(ctx context.Context, userToken entity.UserToken) error
	UpdateUserToken(ctx context.Context, userToken entity.UserToken) error
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}
