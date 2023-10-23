package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"image-gallery/internal/user"
	"time"
)

//TODO: сохранить в папку internal cache/user.go

type User interface {
	Get(ctx context.Context, key string) (*user.User, error)
	Set(ctx context.Context, key string, value *user.User) error
}

type UserCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewUserCache(redisCli *redis.Client) User {
	return &UserCache{redisCli: redisCli}
}

func (c *UserCache) Get(ctx context.Context, key string) (*user.User, error) {
	val := c.redisCli.Get(ctx, key).Val()
	var us *user.User

	err := json.Unmarshal([]byte(val), us)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (c *UserCache) Set(ctx context.Context, key string, value *user.User) error {
	userJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisCli.Set(ctx, key, string(userJson), c.Expiration).Err()
}
