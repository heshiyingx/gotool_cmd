func (m *defaultModel) {{.upperStartCamelObject}}FindOneBy{{.uniqueCombineNameCamel}}(ctx context.Context, {{.uniqueSourceNameAndType}}) (*{{.upperStartCamelObject}}, error) {
        {{.uniqueCacheKeyExpression}}

		var {{.upperStartCamelPrimaryKey}} {{.pkNameType}}
        err := m.db.QueryByCtx(ctx, &{{.upperStartCamelPrimaryKey}},{{.uniqueCacheKeyName}},  func(ctx context.Context, p any, db *gorm.DB) error {
            return db.Model(&{{.upperStartCamelObject}}{}).Select("{{.pkNameWrap}}").Where("{{.originalField}}", {{.lowerStartCamelField}}).Take(p).Error
        })
        if err != nil {
            return nil, err
        }

        {{.pkKeyExpression}}
        var resp {{.upperStartCamelObject}}
        err = m.db.QueryByCtx(ctx, &resp,  {{.pkCacheKeyName}}, func(ctx context.Context, r any, db *gorm.DB) error {
            return db.Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}}= ?", {{.upperStartCamelPrimaryKey}}).Take(r).Error
        })
        return &resp, err

}

func (m *defaultModel) {{.upperStartCamelObject}}UpdateOneBy{{.uniqueCombineNameCamel}}(ctx context.Context, {{.uniqueSourceNameAndType}},updateObj *{{.upperStartCamelObject}},delCacheKeys []string,fields ...string) (int64, error) {
    if updateObj==nil{
		return 0,nil
    }
        data,err := m.{{.upperStartCamelObject}}FindOneBy{{.uniqueCombineNameCamel}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
            return 0, err
		}
		{{.allCacheKeyExpression}}
        delCacheAllKeys := make([]string, 0, {{.allCacheKeyCount}}+len(delCacheKeys))
        delCacheAllKeys = append(delCacheAllKeys, {{.allCacheKeyNames}})
        if len(delCacheKeys) > 0 { delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...) }

        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            upTx := db.Model(&{{.upperStartCamelObject}}{}).Where("{{.pkNameWrap}}", {{.lowerStartCamelField}})
            if len(fields) > 0 {
                upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
            }else {
                upTx = upTx.Save(updateObj)
            }
            return upTx.RowsAffected,upTx.Error
        },delCacheAllKeys...)

}


func (m *defaultModel) {{.upperStartCamelObject}}DeleteOneBy{{.uniqueCombineNameCamel}}(ctx context.Context, {{.uniqueSourceNameAndType}},delCacheKeys ...string) (int64, error) {

        data,err := m.{{.upperStartCamelObject}}FindOneBy{{.uniqueCombineNameCamel}}(ctx,{{.lowerStartCamelField}})
        if err != nil {
        return 0, err
        }
        {{.allCacheKeyExpression}}

        delCacheAllKeys := make([]string, 0, {{.allCacheKeyCount}}+len(delCacheKeys))
        delCacheAllKeys = append(delCacheAllKeys, {{.allCacheKeyNames}})
        if len(delCacheKeys) > 0 { delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...) }

        return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
            delTx := db.Where("{{.pkNameWrap}}", {{.lowerStartCamelField}}).Delete(&{{.upperStartCamelObject}}{})
            return delTx.RowsAffected,delTx.Error
        },delCacheAllKeys...)

    }