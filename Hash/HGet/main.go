package main

import (
	"context"
	"log"

	"go-redis-demo/pkg/redis"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	err := rdb.HSet(ctx, "key", "field", "value").Err()
	if err != nil {
		panic(err)
	}

	res, err := rdb.HGet(ctx, "key", "field").Result()
	if err != nil {
		panic(err)
	}
	log.Println("result: ", res)
}
