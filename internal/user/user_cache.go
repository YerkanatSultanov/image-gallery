package user

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

//TODO: сохранить в папку internal cache/user_cache.go

type UserC interface {
	Get(ctx context.Context, key string) (*User, error)
	Set(ctx context.Context, key string, value *User) error
}

type UserCache struct {
	Expiration time.Duration
	redisCli   *redis.Client
}

func NewUserCache(redisCli *redis.Client) UserC {
	return &UserCache{redisCli: redisCli}
}

func (c *UserCache) Get(ctx context.Context, key string) (*User, error) {
	val := c.redisCli.Get(ctx, key).Val()
	var us *User

	err := json.Unmarshal([]byte(val), us)
	if err != nil {
		return nil, err
	}

	return us, nil
}

func (c *UserCache) Set(ctx context.Context, key string, value *User) error {
	userJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.redisCli.Set(ctx, key, string(userJson), c.Expiration).Err()
}
