package redis

import (
	"GIN/configs"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func ConnectRedis() {
	config := configs.GetConfig()
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if _, err := Client.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Không thể kết nối tới Redis: %v", err)
	}
	log.Println("Connected to Redis:", config.Redis.Addr)

}
