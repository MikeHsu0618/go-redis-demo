package main

import (
	"context"
	"go-redis-demo/pkg/redis"
	"log"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	err := rdb.HSet(ctx, "key", "field", "1", "field2", "100").Err()
	if err != nil {
		panic(err)
	}

	keys, err := rdb.HKeys(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	for i, key := range keys {
		log.Println("index:", i, "key:", key)
	}
}
