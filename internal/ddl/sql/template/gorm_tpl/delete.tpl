{{if .withCache}}
func (m *defaultModel) {{.upperStartCamelObject}}DeleteBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyType}},delCacheKeys ...string) (int64, error) {

        {{if .containsIndexCache}}
            data, err := m.{{.upperStartCamelObject}}FindBy{{.titlePrimaryKey}}(ctx, {{.lowerStartCamelPrimaryKey}})
            if err != nil {
                return 0,err
            }
        {{end}}
        {{.allCacheKeyExpressStr}}

        delCacheAllKeys := make([]string, 0, {{.allCacheKeyCount}}+len(delCacheKeys))
        
		delCacheAllKeys = append(delCacheAllKeys,{{.allCacheKeyNameStr}})
		if len(delCacheKeys) > 0 { delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...) }


        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            res := db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
            return res.RowsAffected, res.Error
		}, delCacheAllKeys...)

}
{{else}}
 func (m *defaultModel) {{.upperStartCamelObject}}DeleteBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.primaryKeyType}},delCacheKeys ...string) (int64, error) {
    return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
    res := db.Where("{{.originalPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete(&{{.upperStartCamelObject}}{})
    return res.RowsAffected, res.Error
    })
}
{{end}}