package main

import (
	"context"
	"log"

	"go-redis-demo/pkg/redis"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	err := rdb.HSet(ctx, "key", "field", "1111").Err()
	if err != nil {
		panic(err)
	}

	err = rdb.HSetNX(ctx, "key", "field", "2222").Err()
	if err != nil {
		panic(err)
	}

	result, err := rdb.HGet(ctx, "key", "field").Result()
	log.Println("result:", result)
}
