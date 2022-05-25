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

	err = rdb.HDel(ctx, "key", "field2").Err()
	if err != nil {
		panic(err)
	}

	results, err := rdb.HMGet(ctx, "key", "field", "field2").Result()

	log.Println(results) // 第二個 field2 刪除後取得會為 nil
}
