package gen

import (
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/collection"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"github.com/heshiyingx/gotool_cmd/internal/ddl/sql/template"
	"sort"
	"strings"
)

func genDeleteByPK(table Table, withCache bool) (string, string, error) {
	uniqueCacheKeyExpressionSet := collection.NewSet[string]()
	uniqueCacheKeyNameSet := collection.NewSet[string]()
	allKeyExpressionSet := collection.NewSet[string]()
	allKeyNameSet := collection.NewSet[string]()
	primaryCacheKeyExpress := table.PrimaryCacheKey.KeyExpression
	primaryCacheKeyName := table.PrimaryCacheKey.KeyLeft

	allKeyExpressionSet.Add(table.PrimaryCacheKey.KeyExpression)
	allKeyNameSet.Add(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		uniqueCacheKeyExpressionSet.Add(key.DataKeyExpression)
		uniqueCacheKeyNameSet.Add(key.KeyLeft)
		allKeyExpressionSet.Add(key.DataKeyExpression)
		allKeyNameSet.Add(key.KeyLeft)
	}
	uniqueCacheKeyNames := uniqueCacheKeyNameSet.Elems()
	sort.Strings(uniqueCacheKeyNames)
	uniqueCacheKeyExpressions := uniqueCacheKeyExpressionSet.Elems()
	sort.Strings(uniqueCacheKeyExpressions)

	allCacheKeyNames := allKeyNameSet.Elems()
	sort.Strings(allCacheKeyNames)
	allCacheKeyExpressions := allKeyExpressionSet.Elems()
	sort.Strings(allCacheKeyExpressions)

	//camel := table.Name.ToCamel()
	text, err := pathext.LoadTemplate(category, deleteTemplateFile, template.Delete)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("delete").
		Parse(text).
		Execute(map[string]any{
			"upperStartCamelObject":     table.Name.ToCamel(), //模型名称
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey, //是否包含唯一索引缓存
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"primaryKeyType":            table.PrimaryKey.Fields[0].DataType,
			"allCacheKeyExpressStr":     strings.Join(allCacheKeyExpressions, "\n"),
			"allCacheKeyNameStr":        strings.Join(allCacheKeyNames, ","),
			"allCacheKeyCount":          len(allCacheKeyNames),
			"primaryCacheKeyName":       primaryCacheKeyName,
			"primaryCacheKeyExpress":    primaryCacheKeyExpress,
			//"pkCacheKey":                table.PrimaryCacheKey.KeyLeft,
			//"uniqueCacheKeys":           keyVars,
			//"uniqueKeysLen":             len(table.UniqueCacheKey),
			//"keysLen":                   len(keys),
			//"cacheKeys":                 append(keyVars, wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source())),
			//"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = pathext.LoadTemplate(category, deleteMethodTemplateFile, template.DeleteMethod)
	if err != nil {
		return "", "", err
	}

	deleteMethodOut, err := util.With("deleteMethod").
		Parse(text).
		Execute(map[string]any{
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Fields[0].Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.Fields[0].DataType,
			"data":                      table,
			"titlePrimaryKey":           table.PrimaryKey.Fields[0].Name.Title(),
			"upperStartCamelObject":     table.Name.ToCamel(),
		})
	if err != nil {
		return "", "", err
	}

	return output.String(), deleteMethodOut.String(), nil
}
