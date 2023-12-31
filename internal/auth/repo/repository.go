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
	CreateUserToken(userToken *entity.UserToken) (*entity.UserToken, error)
	UpdateUserToken(userToken *entity.UserToken) (*entity.UserToken, error)
	GetUserTokenByUserID(userId int) (*entity.UserToken, error)
	IsTokenPresentInDatabase(tokenString string) (bool, error)
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}
