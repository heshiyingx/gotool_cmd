func (m *default{{.upperStartCamelObject}}Model) FindBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
		var resp {{.upperStartCamelObject}}
		err := m.db.QueryOneByPKCtx(ctx, &resp, {{.cacheKeyVariable}}, func(ctx context.Context, r any, db *gorm.DB) error {
			return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}}).Take(r).Error
		})
		return &resp,err
	{{else}}
	var resp {{.upperStartCamelObject}}
	err := m.db.QueryCtx(ctx, &resp, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}}).Take(r).Error
	})
	return &resp, err
	{{end}}
}
