package nouique

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	cacheChickenIdPrefix = "cache:chicken:id:"
)

type (
	chickenModel interface {
		ChickenInsert(ctx context.Context, data *Chicken, delCacheKeys ...string) (int64, error)
		ChickenFindById(ctx context.Context, id int64) (*Chicken, error)
		ChickenUpdateById(ctx context.Context, id int64, updateObj *Chicken, delCacheKeys []string, fields ...string) (int64, error)
		ChickenDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error)
	}

	Chicken struct {
		Id             int64          `db:"id" gorm:"column:id" json:"id,omitempty"`                                        // 主键ID
		UserId         int64          `db:"user_id" gorm:"column:user_id" json:"user_id,omitempty"`                         // 用户ID
		FeedSlotValue  int64          `db:"feed_slot_value" gorm:"column:feed_slot_value" json:"feed_slot_value,omitempty"` // 饲料槽中的饲料数量
		Name           string         `db:"name" gorm:"column:name" json:"name,omitempty"`                                  // 小鸡名字
		FeedNum        int64          `db:"feed_num" gorm:"column:feed_num" json:"feed_num,omitempty"`                      // 喂食次数
		OpTime         int64          `db:"op_time" gorm:"column:op_time" json:"op_time,omitempty"`                         // 上次操作时间
		Stage          int8           `db:"stage" gorm:"column:stage" json:"stage,omitempty"`                               // 所处阶段，1:孵化期，2.成长期，3.下蛋期
		IsDie          int8           `db:"is_die" gorm:"column:is_die" json:"is_die,omitempty"`                            // 是否死亡,0:正常，1:已死亡
		NutritionValue int64          `db:"nutrition_value" gorm:"column:nutrition_value" json:"nutrition_value,omitempty"` // 营养值
		FeedTime       int64          `db:"feed_time" gorm:"column:feed_time" json:"feed_time,omitempty"`                   // 最近一次喂饲料的时间(时间戳秒)
		Process        int64          `db:"process" gorm:"column:process" json:"process,omitempty"`                         // 进度
		CreatedAt      *time.Time     `db:"created_at" gorm:"column:created_at" json:"created_at,omitempty"`                // 创建时间
		UpdatedAt      *time.Time     `db:"updated_at" gorm:"column:updated_at" json:"updated_at,omitempty"`                // 更新时间
		Config         sql.NullString `db:"config" gorm:"column:config" json:"config,omitempty"`                            // 喂养配置快照
	}
)

func (m *defaultModel) ChickenInsert(ctx context.Context, data *Chicken, delCacheKeys ...string) (int64, error) {

	afterDel := true
	delCacheAllKeys := make([]string, 0, 1+len(delCacheKeys))
	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	if data.Id != 0 {
		afterDel = false
		chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, data.Id)
		delCacheAllKeys = append(delCacheAllKeys, chickenIdKey)
	}

	result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.WithContext(ctx).Model(&Chicken{}).Create(data)
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

	if err != nil {
		return 0, err
	}

	if afterDel {
		chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, data.Id)
		err = m.db.DelCacheKeysAndDelay(ctx, chickenIdKey)
		if err != nil {
			return 0, err
		}
	}

	return result, err

}

func (m *defaultModel) ChickenFindById(ctx context.Context, id int64) (*Chicken, error) {
	chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, id)
	var resp Chicken
	err := m.db.QueryByCtx(ctx, &resp, chickenIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.WithContext(ctx).Model(&Chicken{}).Where("`id`=?", id).Take(r).Error
	})
	return &resp, err

}

func (m *defaultModel) ChickenUpdateById(ctx context.Context, id int64, updateObj *Chicken, delCacheKeys []string, fields ...string) (int64, error) {
	if updateObj == nil {
		return 0, nil
	}
	delCacheAllKeys := make([]string, 0, 1+len(delCacheKeys))
	data, err := m.ChickenFindById(ctx, id)
	if err != nil {
		return 0, err
	}

	chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, data.Id)
	delCacheAllKeys = append(delCacheAllKeys, chickenIdKey)
	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&Chicken{}).Where("`id`=?", id)
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		} else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected, upTx.Error
	}, delCacheAllKeys...)

}

func (m *defaultModel) ChickenDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error) {

	chickenIdKey := fmt.Sprintf("%s%v", cacheChickenIdPrefix, id)
	delCacheAllKeys := make([]string, 0, 1+len(delCacheKeys))
	delCacheAllKeys = append(delCacheAllKeys, chickenIdKey)
	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.WithContext(ctx).Where("id = ?", id).Delete(&Chicken{})
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

}

func (Chicken) TableName() string {
	return "chicken"
}
