package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
)

type defaultModel struct {
	db *gormdb.CacheGormDB
}

type DBModel interface {
	FeedStoreDBInterface
	FeedStoreHistoryDBInterface
	NutritionStoreDBInterface
	NutritionStoreHistoryDBInterface
	ChickenDBInterface
}

func NewDBModel(config gormdb.Config) DBModel {

	cacheGormDB := gormdb.MustNewCacheGormDB(config)
	return &defaultModel{
		db: cacheGormDB,
	}

}
