package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"image-gallery/internal/auth/entity"
	"time"
)

func (r *repository) CreateUserToken(userToken *entity.UserToken) (*entity.UserToken, error) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var lastInsertId int
	query := "INSERT INTO user_token(token, refresh_token, user_id) VALUES ($1, $2, $3) returning id"
	err := r.db.QueryRowContext(c, query, userToken.Token, userToken.RefreshToken, userToken.UserId).Scan(&lastInsertId)
	if err != nil {
		return &entity.UserToken{}, err
	}

	userToken.Id = lastInsertId
	return userToken, nil
}

func (r *repository) UpdateUserToken(userToken entity.UserToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	q := `UPDATE user_token SET token = $1, refresh_token = $2 WHERE user_id = $3;
`
	query, args, err := sqlx.In(
		q,
		userToken.Token,
		userToken.RefreshToken,
		userToken.UserId,
	)

	if err != nil {
		return fmt.Errorf("query bake failed: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db exec query failed: %w", err)
	}

	return nil
}

func (r *repository) GetUserTokenByUserID(userId int) (*entity.UserToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `Select id, token, refresh_token, user_id from user_token where user_id = $1`

	u := entity.UserToken{}

	err := r.db.QueryRowContext(ctx, query, userId).Scan(&u.Id, &u.Token, &u.RefreshToken, &u.UserId)
	if err != nil {
		return nil, fmt.Errorf("db exec query failed: %s", err)
	}

	return &u, nil
}
func (r *repository) IsTokenPresentInDatabase(tokenString string) (bool, error) {
	query := "SELECT COUNT(*) FROM user_token WHERE token = $1"
	var count int

	row := r.db.QueryRowContext(context.Background(), query, tokenString)

	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check token in the database: %w", err)
	}

	return count > 0, nil
}
