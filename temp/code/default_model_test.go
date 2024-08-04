package code

import (
	"context"
	"fmt"
	"github.com/heshiyingx/gotool/dbext/gormdb/gormlogs"
	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func TestNewDBModel(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "root",
	})
	model := NewDBModel(gormdb.Config{
		DSN:    "root:root@tcp(127.0.0.1:3306)/chicken_server?charset=utf8mb4&parseTime=True&loc=Local",
		DBType: "MYSQL",
		GormConfig: gorm.Config{
			Logger: gormlogs.Default.LogMode(4),
		},
		Rdb:               client,
		NotFoundExpireSec: 100,
		CacheExpireSec:    100,
		RandSec:           100,
		PreFunc:           nil,
	})
	wg := sync.WaitGroup{}
	wg.Add(300)
	for i := 0; i < 300; i++ {

		go func() {
			defer wg.Done()
			store, err := model.FeedStoreNoCacheNolimitFindOneByUserId(context.Background(), 13)
			//store, err := model.FeedStoreNoCacheFindOneByUserId(context.Background(), 13)
			//store, err := model.FeedStoreFindOneByUserId(context.Background(), 13)
			if err != nil {
				return
			}
			fmt.Println(store)
		}()
	}
	wg.Wait()
}
