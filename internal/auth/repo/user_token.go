package repo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"image-gallery/internal/auth/entity"
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
	q := `
UPDATE user_token SET token = ?, refresh_token = ? WHERE user_id = ?;
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
