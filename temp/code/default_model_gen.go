package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
)

type defaultModel struct {
	db *gormdb.CacheGormDB[int64]
}

type DBModel interface {
	feedStoreModel
	feedStoreHistoryModel
	nutritionStoreModel
	nutritionStoreHistoryModel
	chickenModel
}

func NewDBModel(config gormdb.Config) DBModel {

	cacheGormDB := gormdb.MustNewCacheGormDB[int64](config)
	return &defaultModel{
		db: cacheGormDB,
	}

}
