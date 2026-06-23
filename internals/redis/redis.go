package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedis(host string, port string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect redis:", err)
	}

	fmt.Println("Connected to Redis")

	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}
}
