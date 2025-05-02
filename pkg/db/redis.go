package db

import (
	"awesomeProject/pkg/configs"
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var redisCli redisClient
var once sync.Once

func initRedis() {
	redisCli.client = redis.NewClient(&redis.Options{
		Addr:     configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port,
		Password: configs.GetConfig().Redis.Password,
		DB:       configs.GetConfig().Redis.DB,
	})
}

func NewCache() Cache {
	once.Do(func() {
		initRedis()
	})
	return &redisCli
}

type redisClient struct {
	client *redis.Client
}

func (rc *redisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

func (rc *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}
