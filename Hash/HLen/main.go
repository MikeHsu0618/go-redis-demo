package main

import (
	"context"
	"go-redis-demo/pkg/redis"
	"log"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	err := rdb.HSet(ctx, "key", "field", "1", "field2", "100", "field3", "300").Err()
	if err != nil {
		panic(err)
	}

	count, err := rdb.HLen(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	log.Println("count:", count)
}
