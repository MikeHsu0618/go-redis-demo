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

	err = rdb.HSet(ctx, "key", "field2", "value2").Err()
	if err != nil {
		panic(err)
	}
	log.Println("Done")
}
