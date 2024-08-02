func (m *defaultModel) {{.upperStartCamelObject}}DeleteBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},delCacheKeys ...string) (int64, error) {
    {{if .withCache}}
        {{if .containsIndexCache}}
            data, err := m.{{.upperStartCamelObject}}FindBy{{.titlePrimaryKey}}(ctx, {{.lowerStartCamelPrimaryKey}})
            if err != nil {
                return 0,err
            }
        {{end}}
        {{.keys}}

        delCacheAllKeys := make([]string, 0, {{.keysLen}}+len(delCacheKeys))
		{{ if gt .uniqueKeysLen 0}}
        // 0
		delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}} {{- range $i,$key :=.uniqueCacheKeys }}, {{$key}} {{- end}})
		{{else}}
        // 1
		delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}})
		{{end}}
		if len(delCacheKeys) > 0 { delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...) }


        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.WithContext(ctx).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
		}, delCacheAllKeys...)
    {{else}}
        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.WithContext(ctx).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
        })
    {{end}}
}
