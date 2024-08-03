{{if .withCache}}
func (m *defaultModel) {{.upperStartCamelObject}}FindBy{{.titlePrimaryKey}}(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{.primaryCacheKeyExpress}}
		var resp {{.upperStartCamelObject}}
		err := m.db.QueryByCtx(ctx, &resp, {{.primaryCacheKeyName}}, func(ctx context.Context, r any, db *gorm.DB) error {
			return db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Where("{{.originalPrimaryKey}}=?", {{.lowerStartCamelPrimaryKey}}).Take(r).Error
		})
		return &resp,err
}
{{else}}
func (m *defaultModel) ChickenFindById(ctx context.Context, id int64) (*Chicken, error) {
	chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, id)
	var resp Chicken
	err := m.db.QueryByCtx(ctx, &resp, chickenIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.WithContext(ctx).Model(&Chicken{}).Where("`id`=?", id).Take(r).Error
	})
	return &resp, err

}
{{end}}