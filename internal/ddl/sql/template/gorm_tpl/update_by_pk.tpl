func (m *defaultModel) {{.upperStartCamelObject}}UpdateBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},updateObj *{{.upperStartCamelObject}},delCacheKeys []string,fields ...string) (int64, error) {
	if updateObj==nil{
		return 0,nil
	}
	{{if .withCache}}
	data,err := m.{{.upperStartCamelObject}}FindBy{{.titlePrimaryKey}}(ctx,{{.lowerStartCamelPrimaryKey}})
	if err != nil {
		return 0, err
	}
	{{.keys}}

	delCacheAllKeys := make([]string, 0, {{.keysLen}}+len(delCacheKeys))
		{{ if gt .uniqueKeysLen 0}}
		delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}},{{- range $i,$key :=.uniqueCacheKeys }},{{$key}}{{- end}})
		{{else}}
		delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}})
		{{end}}
	if len(delCacheKeys) > 0 { delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...) }

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}})
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		}else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected,upTx.Error
	},delCacheAllKeys...)


	{{else}}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}})
	if len(fields) > 0 {
		upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
	} else {
		upTx = upTx.Save(updateObj)
	}
		return upTx.RowsAffected, upTx.Error
	})

	{{end}}
}
