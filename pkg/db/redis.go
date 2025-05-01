package db

import (
	"awesomeProject/pkg/configs"
	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port,
		Password: configs.GetConfig().Redis.Password,
		DB:       configs.GetConfig().Redis.DB,
	})
}
