package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"image-gallery/internal/gallery/entity"
	"time"
)

type User interface {
	Get(ctx context.Context, key string) ([]*entity.PhotoResponse, error)
	Set(ctx context.Context, key string, value []*entity.PhotoResponse) error
}

type UserCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewUserCache(redisCli *redis.Client, expiration time.Duration) User {
	return &UserCache{
		redisCli:   redisCli,
		Expiration: expiration,
	}
}

func (u *UserCache) Get(ctx context.Context, key string) ([]*entity.PhotoResponse, error) {
	value := u.redisCli.Get(ctx, key).Val()

	if value == "" {
		return nil, nil
	}

	var user []*entity.PhotoResponse

	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserCache) Set(ctx context.Context, key string, value []*entity.PhotoResponse) error {
	userJson, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return u.redisCli.Set(ctx, key, string(userJson), u.Expiration).Err()
}
