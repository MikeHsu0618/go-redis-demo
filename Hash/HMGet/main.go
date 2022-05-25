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

	results, err := rdb.HMGet(ctx, "key", "field", "field3").Result()
	if err != nil {
		panic(err)
	}

	// 取出來的 value 會是一個 []interface
	for _, result := range results {
		log.Println(result.(string))
	}
}
