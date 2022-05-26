package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-redis-demo/pkg/postgres"
	redispkg "go-redis-demo/pkg/redis"
	"log"
	"strings"
	"sync"
	"time"
)

type Betfair struct {
	EventId string `json:"event_id"`
}

var ctx = context.Background()
var rdb = redispkg.NewRedisClient()

const MmUpdateTimeKey = "mm_update_time"
const TimeFormat = time.RFC3339Nano

var events = getAllDoneEvents() // 取得 betfair 所有已完成的 event_id

func main() {
	wg := new(sync.WaitGroup)
	keyChan := make(chan string, 100)

	// 產生 redis 假資料
	setFakeMMUpdateTime()

	// 取得 hash{ mm_update_time } 所有的 Key 並且丟進 channel 中
	keys := getMMUpdateTimeKeys()
	wg.Add(len(keys))

	// 開啟 transaction pipeline
	pipe := rdb.TxPipeline()

	// 開啟 10 個 worker -> 取得一個 key 開始跟 event_id 比對, 如果在其中就 HDel Key
	for i := 0; i < 10; i++ {
		go worker(keyChan, wg, pipe)
	}

	// 插入 keyChan
	for _, key := range keys {
		keyChan <- key
	}

	wg.Wait()
	// transaction pipeline commit
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("pipeline error")
		return
	}
	log.Println("done")
}

func worker(keyChan <-chan string, wg *sync.WaitGroup, pipe redis.Pipeliner) {
	for key := range keyChan {
		eventId := strings.Split(key, "_")[0]
		for _, event := range events {
			if eventId != event.EventId {
				continue
			}
			log.Printf("刪除 key : %v, id: %v", key, event.EventId)
			pipe.HDel(ctx, MmUpdateTimeKey, key)
		}
		wg.Done()
	}
}

func getAllDoneEvents() []Betfair {
	db := postgres.NewPostgresClient()
	var result []Betfair
	db.Table("betfair").Select("event_id").Where("settle_dt != ?", "0001-01-01 00:00:00").Find(&result)

	return result
}

func setFakeMMUpdateTime() {
	tstring := time.Now().Format(TimeFormat)
	fieldName := "31104929_1.198617688_44720863_b"
	pipe := rdb.TxPipeline()
	for i := 0; i < 10000; i++ {
		err := pipe.HSet(ctx, MmUpdateTimeKey, fmt.Sprintf("%v%v", fieldName, string(i)), tstring).Err()
		if err != nil {
			log.Println(err)
		}
	}
	pipe.Exec(ctx)
}

func getMMUpdateTimeKeys() []string {
	res, err := rdb.HKeys(ctx, MmUpdateTimeKey).Result()
	if err != nil {
		log.Println(err)
	}

	return res
}
