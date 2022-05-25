package main

import (
	"context"
	"go-redis-demo/pkg/redis"
	"log"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	err := rdb.HSet(ctx, "key", "field", "value", "field2", "value").Err()
	if err != nil {
		panic(err)
	}

	results, err := rdb.HGetAll(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	// data 是個 map[field]value 類型
	for field, value := range results {
		log.Println("field: ", field, "value : ", value)
	}
}
