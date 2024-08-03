// Code generated by gotool. DO NOT EDIT.
// If you find any bugs, please contact heshiyingx@126.com.
// Your help is greatly appreciated.

package code

import (
	"context"

	"fmt"
	"strings"
	"time"

	// "github.com/heshiyingx/gotool/dbext/gormdb/v2"
	"gorm.io/gorm"
)

var (
	cacheNutritionStoreHistoryIdPrefix         = "cache:nutritionStoreHistory:id:"
	cacheNutritionStoreHistoryUserIdOpIdPrefix = "cache:nutritionStoreHistory:userId:opId:"
)

type (
	nutritionStoreHistoryModel interface {
		NutritionStoreHistoryInsert(ctx context.Context, data *NutritionStoreHistory, delCacheKeys ...string) (int64, error)
		NutritionStoreHistoryFindById(ctx context.Context, id int64) (*NutritionStoreHistory, error)
		NutritionStoreHistoryUpdateById(ctx context.Context, id int64, updateObj *NutritionStoreHistory, delCacheKeys []string, fields ...string) (int64, error)
		NutritionStoreHistoryDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error)
	}

	NutritionStoreHistory struct {
		Id        int64      `db:"id" gorm:"column:id" json:"id,omitempty"`                         // 主键
		UserId    int64      `db:"user_id" gorm:"column:user_id" json:"user_id,omitempty"`          // 用户ID
		ChickenId int64      `db:"chicken_id" gorm:"column:chicken_id" json:"chicken_id,omitempty"` // 小鸡ID
		OpType    int64      `db:"op_type" gorm:"column:op_type" json:"op_type,omitempty"`          // 操作类型,1:增加，2：减少
		Value     int64      `db:"value" gorm:"column:value" json:"value,omitempty"`                // 改变的值
		OpId      string     `db:"op_id" gorm:"column:op_id" json:"op_id,omitempty"`                // 操作id，用于去重
		Comment   string     `db:"comment" gorm:"column:comment" json:"comment,omitempty"`          // 备注说明
		CreatedAt *time.Time `db:"created_at" gorm:"column:created_at" json:"created_at,omitempty"` // 创建时间
		UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at" json:"updated_at,omitempty"` // 更新时间
	}
)

func (m *defaultModel) NutritionStoreHistoryInsert(ctx context.Context, data *NutritionStoreHistory, delCacheKeys ...string) (int64, error) {
	nutritionStoreHistoryUserIdOpIdKey := fmt.Sprintf("%s%v:%v", cacheNutritionStoreHistoryUserIdOpIdPrefix, data.UserId, data.OpId)
	afterDel := true

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))
	delCacheAllKeys = append(delCacheAllKeys, nutritionStoreHistoryUserIdOpIdKey)

	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	if data.Id != 0 {
		afterDel = false
		nutritionStoreHistoryIdKey := fmt.Sprintf("%s%v", cacheNutritionStoreHistoryIdPrefix, data.Id)
		delCacheAllKeys = append(delCacheAllKeys, nutritionStoreHistoryIdKey)
	}

	result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.WithContext(ctx).Model(&NutritionStoreHistory{}).Create(data)
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

	if err != nil {
		return 0, err
	}

	if afterDel {
		nutritionStoreHistoryIdKey := fmt.Sprintf("%s%v", cacheNutritionStoreHistoryIdPrefix, data.Id)
		err = m.db.DelCacheKeysAndDelay(ctx, nutritionStoreHistoryIdKey)
		if err != nil {
			return 0, err
		}
	}
	return result, err

}

func (m *defaultModel) NutritionStoreHistoryFindById(ctx context.Context, id int64) (*NutritionStoreHistory, error) {
	nutritionStoreHistoryIdKey := fmt.Sprintf("%s%v", cacheNutritionStoreHistoryIdPrefix, id)
	var resp NutritionStoreHistory
	err := m.db.QueryByCtx(ctx, &resp, nutritionStoreHistoryIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.WithContext(ctx).Model(&NutritionStoreHistory{}).Where("`id`=?", id).Take(r).Error
	})
	return &resp, err
}

func (m *defaultModel) NutritionStoreHistoryUpdateById(ctx context.Context, id int64, updateObj *NutritionStoreHistory, delCacheKeys []string, fields ...string) (int64, error) {
	if updateObj == nil {
		return 0, nil
	}

	data, err := m.NutritionStoreHistoryFindById(ctx, id)
	if err != nil {
		return 0, err
	}
	nutritionStoreHistoryIdKey := fmt.Sprintf("%s%v", cacheNutritionStoreHistoryIdPrefix, data.Id)
	nutritionStoreHistoryUserIdOpIdKey := fmt.Sprintf("%s%v:%v", cacheNutritionStoreHistoryUserIdOpIdPrefix, data.UserId, data.OpId)

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))

	delCacheAllKeys = append(delCacheAllKeys, nutritionStoreHistoryIdKey)

	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&NutritionStoreHistory{}).Where("`id`=?", id)
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		} else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected, upTx.Error
	}, delCacheAllKeys...)

}

func (m *defaultModel) NutritionStoreHistoryDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error) {

	data, err := m.NutritionStoreHistoryFindById(ctx, id)
	if err != nil {
		return 0, err
	}

	nutritionStoreHistoryIdKey := fmt.Sprintf("%s%v", cacheNutritionStoreHistoryIdPrefix, id)
	nutritionStoreHistoryUserIdOpIdKey := fmt.Sprintf("%s%v:%v", cacheNutritionStoreHistoryUserIdOpIdPrefix, data.UserId, data.OpId)

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))

	delCacheAllKeys = append(delCacheAllKeys, nutritionStoreHistoryIdKey, nutritionStoreHistoryUserIdOpIdKey)
	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.Where("`id` = ?", id).Delete(&NutritionStoreHistory{})
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

}

func (NutritionStoreHistory) TableName() string {
	return "nutrition_store_history"
}
