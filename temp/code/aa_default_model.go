package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
)

type defaultModel struct {
	db *gormdb.CacheGormDB[int64]
}

type DBModel interface {
	FeedStoreDBInterface
	FeedStoreHistoryDBInterface
	NutritionStoreDBInterface
	NutritionStoreHistoryDBInterface
	ChickenDBInterface
}

func NewDBModel(config gormdb.Config) DBModel {

	cacheGormDB := gormdb.MustNewCacheGormDB[int64](config)
	return &defaultModel{
		db: cacheGormDB,
	}

}
