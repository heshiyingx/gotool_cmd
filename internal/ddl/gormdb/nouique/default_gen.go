package nouique

import (
	"github.com/heshiyingx/gotool_cmd/internal/ddl/gormdb"
)

type defaultModel struct {
	db *gormdb.CacheGormDB[int64]
}

type DBModel interface {
	ChickenInterface
	UserChickenInterface
}

func NewDBModel(config gormdb.Config) DBModel {

	cacheGormDB := gormdb.MustNewCacheGormDB[int64](config)
	return &defaultModel{
		db: cacheGormDB,
	}

}
