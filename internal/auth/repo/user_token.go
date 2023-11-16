package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"image-gallery/internal/auth/entity"
	"time"
)

func (r *repository) CreateUserToken(ctx context.Context, userToken entity.UserToken) error {
	q := `insert into user_token(token, refresh_token, user_id) values ($1, $2, $3) returning id`
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

func (r *repository) UpdateUserToken(ctx context.Context, userToken entity.UserToken) error {
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
