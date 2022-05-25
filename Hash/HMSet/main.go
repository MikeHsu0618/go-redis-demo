package main

import (
	"context"
	"go-redis-demo/pkg/redis"
	"log"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewRedisClient()

	data := make(map[string]string)

	data["field1"] = "100"
	data["field2"] = "value"
	data["field3"] = "200"

	err := rdb.HMSet(ctx, "key", data).Err()
	if err != nil {
		panic("error")
	}

	result, err1 := rdb.HGetAll(ctx, "key").Result()
	if err1 != nil {
		panic(err1)
	}

	log.Println(result)
}
