package bootstrap

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	redisClient *redis.Client
	once        sync.Once
	ctx         = context.Background()
)

// RedisClient ...
type RedisClient struct {
	address  string
	password string
	db       int
}

type IRedisClient interface {
	SetRedis(ctx context.Context, key string, value string, ttl int) error
	GetRedis(ctx context.Context, key string) error
}

// NewRedisClient is an interface for managing tokens
// address string, password string, db int
func NewRedisClient(env *Env) *RedisClient {
	return &RedisClient{
		address:  fmt.Sprintf("%s:%s", env.CacheRedisHost, env.CacheRedisPort),
		password: env.CacheRedisPassword,
		db:       env.CacheRedisDB,
	}
}

func (r *RedisClient) RedisConn() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     r.address,
			Password: r.password,
			DB:       r.db, // use the last database
		})
	})
	return redisClient
}

func (r *RedisClient) SetRedis(ctx context.Context, key string, value string, ttl int) error {
	var err error
	if ttl == 0 {
		err = r.RedisConn().Set(ctx, key, value, 0).Err()
	} else {
		err = r.RedisConn().Set(ctx, key, value, time.Duration(ttl)*time.Minute).Err()
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetRedis(ctx context.Context, key string) (interface{}, error) {
	res, err := r.RedisConn().Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *RedisClient) DelRedis(ctx context.Context, key string) error {
	err := r.RedisConn().Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
