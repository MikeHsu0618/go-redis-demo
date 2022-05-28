package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-redis-demo/pkg/postgres"
	redispkg "go-redis-demo/pkg/redis"
	"gorm.io/gorm"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Betfair struct {
	EventId string `json:"event_id"`
}

type MMUpdateTime struct {
	Field      string
	UpdateTime string
}

type MMUpdateTimeRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

const (
	MmUpdateTimeKey = "mm_update_time"
	TimeFormat      = time.RFC3339Nano
	ExpireTime      = "72h"
)

var (
	ctx     = context.Background()
	wg      = new(sync.WaitGroup)
	jobChan = make(chan MMUpdateTime, 1000)
)

func main() {

	db := postgres.NewPostgresClient()
	rdb := redispkg.NewRedisClient()
	repo := NewMMUpdateTimeRepo(db, rdb)

	// 產生 redis 假資料
	repo.setFakeMMUpdateTime()
	// 取得 betfair 所有已完成的 event_id
	events := repo.getAllDoneEvents()
	// 取得 hash{ mm_update_time } 所有的 Key 並且丟進 channel 中
	results := repo.getAllMMUpdateTime()
	wg.Add(len(results))

	pipe := repo.rdb.TxPipeline()
	// 開啟 10 個 worker -> 取得一個 field 開始跟 event_id 比對, 如果在其中就 HDel field
	for i := 0; i < runtime.NumCPU(); i++ {
		go repo.worker(jobChan, wg, pipe, events)
	}
	// 插入 jobChan
	for field, value := range results {
		jobChan <- MMUpdateTime{field, value}
	}
	close(jobChan)
	wg.Wait()
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("pipeline error")
		return
	}
	log.Println("done")
}

func NewMMUpdateTimeRepo(db *gorm.DB, rdb *redis.Client) *MMUpdateTimeRepository {
	return &MMUpdateTimeRepository{
		db:  db,
		rdb: rdb,
	}
}

func (repo *MMUpdateTimeRepository) worker(jobChan <-chan MMUpdateTime, wg *sync.WaitGroup, pipe redis.Pipeliner, events []Betfair) {
	for mmUpdateTime := range jobChan {
		eventId := strings.Split(mmUpdateTime.Field, "_")[0]
		for _, event := range events {
			if eventId != event.EventId {
				continue
			}
			t, _ := time.Parse(TimeFormat, mmUpdateTime.UpdateTime)
			expireTime, _ := time.ParseDuration(ExpireTime)
			if t.Add(expireTime).After(time.Now()) {
				continue
			}
			log.Printf("刪除 key : %v, id: %v", mmUpdateTime.Field, event.EventId)
			repo.deleteMMUpdateTime(ctx, pipe, mmUpdateTime.Field)
		}
		wg.Done()
	}
}

func (repo *MMUpdateTimeRepository) deleteMMUpdateTime(ctx context.Context, pipe redis.Pipeliner, field string) {
	pipe.HDel(ctx, MmUpdateTimeKey, field)
}

func (repo *MMUpdateTimeRepository) getAllDoneEvents() []Betfair {
	var result []Betfair
	repo.db.Table("betfair").Select("event_id").Where("settle_dt != ?", "0001-01-01 00:00:00").Find(&result)
	return result
}

func (repo *MMUpdateTimeRepository) setFakeMMUpdateTime() {
	tstring := time.Now().Format(TimeFormat)
	fieldName := "31104929_1.198617688_44720863_b"
	pipe := repo.rdb.TxPipeline()
	for i := 0; i < 10000; i++ {
		err := pipe.HSet(ctx, MmUpdateTimeKey, fmt.Sprintf("%v%v", fieldName, strconv.Itoa(i)), tstring).Err()
		if err != nil {
			log.Println(err)
		}
	}
	pipe.Exec(ctx)
}

func (repo *MMUpdateTimeRepository) getAllMMUpdateTime() map[string]string {
	res, err := repo.rdb.HGetAll(ctx, MmUpdateTimeKey).Result()
	if err != nil {
		log.Println(err)
	}
	return res
}
