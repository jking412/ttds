package db

import (
	"awesomeProject/pkg/configs"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     configs.GetConfig().Redis.Host + ":" + configs.GetConfig().Redis.Port,
		Password: configs.GetConfig().Redis.Password,
		DB:       configs.GetConfig().Redis.DB,
	})
}
