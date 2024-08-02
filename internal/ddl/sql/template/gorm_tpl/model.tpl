package {{.pkg}}
{{if .withCache}}
import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)
{{else}}
import (
	 "github.com/heshiyingx/gotool/dbext/gormdb"
)
{{end}}


type (
	// {{.upperStartCamelObject}}DBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}DBInterface interface {
		{{.lowerStartCamelObject}}Model
	}


)



{{if not .withCache}}

{{end}}
