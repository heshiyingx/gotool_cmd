package nouique

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	cacheUserChickenIdPrefix     = "cache:userChicken:id:"
	cacheUserChickenUserIdPrefix = "cache:userChicken:userId:"
)

type (
	userChickenModel interface {
		UserChickenInsert(ctx context.Context, data *UserChicken, delCacheKeys ...string) (int64, error)
		UserChickenFindById(ctx context.Context, id int64) (*UserChicken, error)
		UserChickenUpdateById(ctx context.Context, id int64, updateObj *UserChicken, delCacheKeys []string, fields ...string) (int64, error)
		UserChickenDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error)
		UserChickenFindOneByUserId(ctx context.Context, userId int64) (*UserChicken, error)
		UserChickenDeleteOneByUserId(ctx context.Context, userId int64, delCacheKeys ...string) (int64, error)
		UserChickenUpdateOneByUserId(ctx context.Context, userId int64, updateObj *UserChicken, delCacheKeys []string, fields ...string) (int64, error)
	}

	UserChicken struct {
		Id          int64      `db:"id" gorm:"column:id" json:"id,omitempty"`
		UserId      int64      `db:"user_id" gorm:"column:user_id" json:"user_id,omitempty"`                   // 用户ID
		EggGenCount int64      `db:"egg_gen_count" gorm:"column:egg_gen_count" json:"egg_gen_count,omitempty"` // 当前生蛋数量
		ReadCount   int64      `db:"read_count" gorm:"column:read_count" json:"read_count,omitempty"`          // 客户端已知晓的生蛋数量
		ChickenId   int64      `db:"chicken_id" gorm:"column:chicken_id" json:"chicken_id,omitempty"`          // 小鸡ID
		CreatedAt   *time.Time `db:"created_at" gorm:"column:created_at" json:"created_at,omitempty"`          // 创建时间
		UpdatedAt   *time.Time `db:"updated_at" gorm:"column:updated_at" json:"updated_at,omitempty"`          // 更新时间
	}
)

func (m *defaultModel) UserChickenInsert(ctx context.Context, data *UserChicken, delCacheKeys ...string) (int64, error) {
	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, data.UserId)
	afterDel := true
	cacheKeys := make([]string, 0, 2)
	cacheKeys = []string{userChickenUserIdKey}
	if len(delCacheKeys) > 0 {
		cacheKeys = append(cacheKeys, delCacheKeys...)
	}

	if data.Id != 0 {
		afterDel = false
		userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, data.Id)
		cacheKeys = append(cacheKeys, userChickenIdKey)
	}

	result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.Model(&UserChicken{}).Create(data)
		return res.RowsAffected, res.Error
	}, cacheKeys...)

	if err != nil {
		return 0, err
	}

	if afterDel {
		userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, data.Id)
		err = m.db.DelCacheKeys(ctx, userChickenIdKey)
		if err != nil {
			return 0, err
		}
	}
	return result, err

}

func (m *defaultModel) UserChickenFindById(ctx context.Context, id int64) (*UserChicken, error) {
	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, id)
	var resp UserChicken
	err := m.db.QueryByCtx(ctx, &resp, userChickenIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.Model(&UserChicken{}).Where("`id`=?", id).Take(r).Error
	})
	return &resp, err

}

func (m *defaultModel) UserChickenUpdateById(ctx context.Context, id int64, updateObj *UserChicken, delCacheKeys []string, fields ...string) (int64, error) {
	if updateObj == nil {
		return 0, nil
	}

	data, err := m.UserChickenFindById(ctx, id)
	if err != nil {
		return 0, err
	}
	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, data.Id)
	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, data.UserId)

	delKeys := []string{userChickenIdKey, userChickenUserIdKey}
	if len(delCacheKeys) > 0 {
		delKeys = append(delKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.Model(&UserChicken{}).Where("`id`=?", id)
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		} else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected, upTx.Error
	}, delKeys...)

}

func (m *defaultModel) UserChickenDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error) {

	data, err := m.UserChickenFindById(ctx, id)
	if err != nil {
		return 0, err
	}

	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, id)
	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, data.UserId)
	delKeys := []string{userChickenIdKey, userChickenUserIdKey}
	if len(delCacheKeys) > 0 {
		delKeys = append(delKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.Where("id = ?", id).Delete(&UserChicken{})
		return res.RowsAffected, res.Error
	}, delKeys...)

}

func (m *defaultModel) UserChickenFindOneByUserId(ctx context.Context, userId int64) (*UserChicken, error) {

	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, userId)

	var Id int64
	err := m.db.QueryByCtx(ctx, &Id, userChickenUserIdKey, func(ctx context.Context, p any, db *gorm.DB) error {
		return db.Model(&UserChicken{}).Select("`id`").Where("`user_id` = ?", userId).Take(p).Error
	})
	if err != nil {
		return nil, err
	}
	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, Id)
	var resp UserChicken
	err = m.db.QueryByCtx(ctx, &resp, userChickenIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.Model(&UserChicken{}).Where("`id` = ?", Id).Take(r).Error
	})
	return &resp, err

}

func (m *defaultModel) UserChickenUpdateOneByUserId(ctx context.Context, userId int64, updateObj *UserChicken, delCacheKeys []string, fields ...string) (int64, error) {
	if updateObj == nil {
		return 0, nil
	}

	data, err := m.UserChickenFindOneByUserId(ctx, userId)
	if err != nil {
		return 0, err
	}
	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, data.Id)
	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, data.UserId)

	delKeys := []string{userChickenIdKey, userChickenUserIdKey}
	if len(delCacheKeys) > 0 {
		delKeys = append(delKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.Model(&UserChicken{}).Where("`id`", userId)
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		} else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected, upTx.Error
	}, delKeys...)

}
func (m *defaultModel) UserChickenDeleteOneByUserId(ctx context.Context, userId int64, delCacheKeys ...string) (int64, error) {

	data, err := m.UserChickenFindOneByUserId(ctx, userId)
	if err != nil {
		return 0, err
	}
	userChickenIdKey := fmt.Sprintf("%s%v", cacheUserChickenIdPrefix, data.Id)
	userChickenUserIdKey := fmt.Sprintf("%s%v", cacheUserChickenUserIdPrefix, data.UserId)

	delKeys := []string{userChickenIdKey, userChickenUserIdKey}
	if len(delCacheKeys) > 0 {
		delKeys = append(delKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		delTx := db.Where("`id`", userId).Delete(&UserChicken{})
		return delTx.RowsAffected, delTx.Error
	}, delKeys...)

}

func (UserChicken) TableName() string {
	return "user_chicken"
}
