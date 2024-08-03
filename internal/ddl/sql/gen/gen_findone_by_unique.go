package gen

import (
	"fmt"
	"github.com/heshiyingx/gotool/util"
	"github.com/heshiyingx/gotool/util/collection"
	"github.com/heshiyingx/gotool/util/pathext"
	stringx "github.com/heshiyingx/gotool/util/stringext"
	"github.com/heshiyingx/gotool_cmd/internal/ddl/sql/template"
	"sort"
	"strings"
)

type findOneCode struct {
	findOneMethod          string
	findOneInterfaceMethod string
	cacheExtra             string
}

func genFindAndUpdateOneByUniqueKey(table Table, withCache bool) (*findOneCode, error) {
	text, err := pathext.LoadTemplate(category, findOneByFieldTemplateFile, template.FindOneByField)
	if err != nil {
		return nil, err
	}

	t := util.With("findOneByField").Parse(text)
	var list []string
	//camelTableName := table.Name.ToCamel()
	for _, key := range table.UniqueCacheKey {
		allKeyExpressionSet := collection.NewSet[string]()
		allKeyNameSet := collection.NewSet[string]()
		allKeyNameSet.Add(table.PrimaryCacheKey.KeyLeft)
		allKeyExpressionSet.Add(table.PrimaryCacheKey.DataKeyExpression)
		for _, key := range table.UniqueCacheKey {
			allKeyExpressionSet.Add(key.DataKeyExpression)
			allKeyNameSet.Add(key.KeyLeft)
		}
		allCacheKeyName := allKeyNameSet.Elems()
		sort.Strings(allCacheKeyName)
		allCacheKeyExpression := allKeyExpressionSet.Elems()
		sort.Strings(allCacheKeyExpression)

		uniqueCacheKeyExpression := key.KeyExpression
		uniqueCacheKeyName := key.KeyLeft
		//uniqueCacheKeyPrefix := key.VarRight
		//primaryKeyNameSource := table.PrimaryKey.Fields[0].NameOriginal
		//primaryKeyNameUpstartCamel := table.PrimaryKey.Fields[0].NameOriginal
		primaryNameSourceWarp := wrapWithRawString(table.PrimaryKey.Fields[0].NameOriginal)
		primaryKeyType := table.PrimaryKey.Fields[0].DataType

		uniqueSourceNameAndType, paramJoinString, originalFieldString := convertJoin(key)

		output, err := t.Execute(map[string]any{
			"upperStartCamelObject":     table.Name.ToCamel(),
			"lowerStartCamelObject":     stringx.From(table.Name.ToCamel()).Untitle(),
			"uniqueCombineNameCamel":    key.FieldNameJoin.Camel().With("").Source(),
			"uniqueSourceNameAndType":   uniqueSourceNameAndType,
			"withCache":                 withCache,
			"uniqueCacheKeyExpression":  uniqueCacheKeyExpression,
			"uniqueCacheKeyName":        uniqueCacheKeyName,
			"pkKeyExpression":           fmt.Sprintf(`%v:=fmt.Sprintf("%s%v", %v, %v)`, table.PrimaryCacheKey.KeyLeft, "%s", "%v", table.PrimaryCacheKey.VarLeft, table.PrimaryKey.Fields[0].Name.ToCamel()),
			"pkNameWrap":                primaryNameSourceWarp,
			"pkNameType":                primaryKeyType,
			"pkCacheKeyName":            table.PrimaryCacheKey.KeyLeft,
			"lowerStartCamelField":      paramJoinString,
			"upperStartCamelPrimaryKey": table.PrimaryKey.Fields[0].Name.ToCamel(),
			"originalField":             originalFieldString,
			"data":                      table,
			"allCacheKeyExpression":     strings.Join(allCacheKeyExpression, "\n"),
			"allCacheKeyNames":          strings.Join(allCacheKeyName, ","),
			"allCacheKeyCount":          len(allCacheKeyName),
		})
		if err != nil {
			return nil, err
		}

		list = append(list, output.String())
	}

	text, err = pathext.LoadTemplate(category, findOneByFieldMethodTemplateFile, template.FindOneByFieldMethod)
	if err != nil {
		return nil, err
	}

	t = util.With("findOneByFieldMethod").Parse(text)
	var listMethod []string
	for _, key := range table.UniqueCacheKey {
		var inJoin, paramJoin Join
		for _, f := range key.Fields {
			param := util.EscapeGolangKeyword(stringx.From(f.Name.ToCamel()).Untitle())
			inJoin = append(inJoin, fmt.Sprintf("%s %s", param, getNotNullType(f.DataType)))
			paramJoin = append(paramJoin, param)
		}

		var in string
		if len(inJoin) > 0 {
			in = inJoin.With(", ").Source()
		}
		output, err := t.Execute(map[string]any{
			"upperStartCamelObject": table.Name.ToCamel(),
			"upperField":            key.FieldNameJoin.Camel().With("").Source(),
			"in":                    in,
			"data":                  table,
			"pkKey":                 fmt.Sprintf(`%v:=fmt.Sprintf("%s%v", %v, %v)`, table.PrimaryCacheKey.KeyLeft, "%s", "%v", table.PrimaryCacheKey.VarLeft, table.PrimaryCacheKey.Fields[0].NameOriginal),
			"pkNameWrap":            wrapWithRawString(table.PrimaryKey.Fields[0].NameOriginal),
			"pkNameType":            table.PrimaryKey.Fields[0].DataType,
			"pkCacheKeyName":        table.PrimaryCacheKey.KeyLeft,
		})
		if err != nil {
			return nil, err
		}

		listMethod = append(listMethod, output.String())
	}

	if withCache {
		text, err := pathext.LoadTemplate(category, findOneByFieldExtraMethodTemplateFile,
			template.FindOneByFieldExtraMethod)
		if err != nil {
			return nil, err
		}

		out, err := util.With("findOneByFieldExtraMethod").Parse(text).Execute(map[string]any{
			"upperStartCamelObject": table.Name.ToCamel(),
			"primaryKeyLeft":        table.PrimaryCacheKey.VarLeft,
			"lowerStartCamelObject": stringx.From(table.Name.ToCamel()).Untitle(),
			"originalPrimaryField":  wrapWithRawString(table.PrimaryKey.Fields[0].Name.Source()),
			"data":                  table,
		})
		if err != nil {
			return nil, err
		}

		return &findOneCode{
			findOneMethod:          strings.Join(list, "\n"),
			findOneInterfaceMethod: strings.Join(listMethod, "\n"),
			cacheExtra:             out.String(),
		}, nil
	}

	return &findOneCode{
		findOneMethod:          strings.Join(list, "\n"),
		findOneInterfaceMethod: strings.Join(listMethod, "\n"),
	}, nil
}

func convertJoin(key Key) (in, paramJoinString, originalFieldString string) {
	var inJoin, paramJoin, argJoin Join
	for _, f := range key.Fields {
		param := util.EscapeGolangKeyword(stringx.From(f.Name.ToCamel()).Untitle())
		inJoin = append(inJoin, fmt.Sprintf("%s %s", param, getNotNullType(f.DataType)))
		paramJoin = append(paramJoin, param)

		//if postgreSql {
		//	argJoin = append(argJoin, fmt.Sprintf("%s = $%d", wrapWithRawString(f.Name.Source()), index+1))
		//} else {
		argJoin = append(argJoin, fmt.Sprintf("%s = ?", wrapWithRawString(f.Name.Source())))
		//}
	}
	if len(inJoin) > 0 {
		in = inJoin.With(", ").Source()
	}

	if len(paramJoin) > 0 {
		paramJoinString = paramJoin.With(",").Source()
	}

	if len(argJoin) > 0 {
		originalFieldString = argJoin.With(" and ").Source()
	}
	return in, paramJoinString, originalFieldString
}
