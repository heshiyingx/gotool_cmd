func (m *defaultModel) {{.upperStartCamelObject}}Insert(ctx context.Context, data *{{.upperStartCamelObject}},delCacheKeys ...string) (int64, error) {
{{if .withCache}}{{.keys}}
		afterDel := true

		delCacheAllKeys := make([]string, 0, {{.keysLen}}+len(delCacheKeys))
		{{ if gt .uniqueKeysLen 0}}delCacheAllKeys = append(delCacheAllKeys, {{- range $i,$key :=.uniqueCacheKeys }}{{ if gt $i 0 }},{{end}}{{$key}}{{- end}}){{end}}

		if len(delCacheKeys) > 0 {
			delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
		}


		{{ if eq .pkDataType "string" }}
			if data.Id != "" {
				afterDel = false
				{{.pkCacheKeyExpression}}
				delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}})
			}
		{{else}}
			if data.Id != 0 {
				afterDel = false
				{{.pkCacheKeyExpression}}
				delCacheAllKeys = append(delCacheAllKeys, {{.pkCacheKey}})
			}
		{{end}}


		result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
			res := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Create(data)
			return res.RowsAffected, res.Error
		}, delCacheAllKeys...)

		if err != nil {
			return 0, err
		}

		if afterDel {
			{{.pkCacheKeyExpression}}
			err = m.db.DelCacheKeysAndDelay(ctx, {{.pkCacheKey}})
			if err != nil {
				return 0, err
			}
		}
		return result, err
{{else}}
		return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
			res := db.WithContext(ctx).Model(&{{.upperStartCamelObject}}{}).Create(data)
			return res.RowsAffected, res.Error
		})
{{end}}
}
