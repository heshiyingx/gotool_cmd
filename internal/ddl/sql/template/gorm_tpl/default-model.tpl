package {{.pkg}}
{{if .withCache}}
import (
	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
)
{{else}}
import (
	 "github.com/heshiyingx/gotool/dbext/gormdb/v2"
)
{{end}}


type defaultModel struct {
	db *gormdb.CacheGormDB[{{.pkType}}]
}

type DBModel interface {
	{{.subModelInterface}}
}

func NewDBModel(config gormdb.Config) DBModel {

	cacheGormDB := gormdb.MustNewCacheGormDB[{{.pkType}}](config)
	return &defaultModel{
		db: cacheGormDB,
	}

}




{{if not .withCache}}

{{end}}
