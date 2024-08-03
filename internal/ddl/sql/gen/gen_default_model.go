package gen

import (
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/pathext"
	"github.com/heshiyingx/gotool_cmd/internal/ddl/sql/parser"
	"github.com/heshiyingx/gotool_cmd/internal/ddl/sql/template"

	goformat "go/format"
	"strings"
)

func genDefaultDBModel(tables []*parser.Table, withCache bool, pkg string) (string, error) {
	pkType := tables[0].PrimaryKey.Fields[0].DataType
	modelInterfaceName := make([]string, 0, len(tables))
	for _, e := range tables {
		modelInterfaceName = append(modelInterfaceName, e.Name.ToCamel()+"DBInterface")
	}
	if withCache {
		text, err := pathext.LoadTemplate(category, defaultDbModelTemplateFile, template.DefaultModel)
		if err != nil {
			return "", err
		}

		buffer, err := util.With("import").Parse(text).Execute(map[string]any{
			"subModelInterface": strings.Join(modelInterfaceName, "\n"),
			"pkType":            pkType,
			"pkg":               pkg,
		})
		if err != nil {
			return "", err
		}
		source, err := goformat.Source(buffer.Bytes())
		if err != nil {
			return "", err
		}
		return string(source), err

	}

	text, err := pathext.LoadTemplate(category, importsWithNoCacheTemplateFile, template.ImportsNoCache)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]any{
		"subModelInterface": strings.Join(modelInterfaceName, "\n"),
		"pkType":            pkType,
		"pkg":               pkg,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
