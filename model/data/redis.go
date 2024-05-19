package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	RedisConfig RedisConfig
	Logger      *zap.Logger
	client      *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (r *RedisClient) StartRedis() {
	r.Logger.Info("Starting Redis client")
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.RedisConfig.Addr,
		Password: r.RedisConfig.Password, // no password set
		DB:       r.RedisConfig.DB,       // use default DB
	})
	r.Logger.Info("Redis client started")

}

func (r *RedisClient) PutInCache(ctx context.Context, key string, value interface{}, ttl time.Duration, sendChan chan<- bool) {
	res := r.client.Set(ctx, key, value, ttl)
	if res.Err() != nil {
		r.Logger.Error("Writing to cache failed")
		sendChan <- false
	} else {
		r.Logger.Info("Written successfully to cache")
		sendChan <- true
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}
