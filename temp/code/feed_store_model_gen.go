// Code generated by gotool. DO NOT EDIT.
// If you find any bugs, please contact heshiyingx@126.com.
// Your help is greatly appreciated.

package code

import (
	"context"

	"fmt"
	"strings"
	"time"

	"github.com/heshiyingx/gotool/dbext/gormdb/v2"
	"gorm.io/gorm"
)

var (
	cacheFeedStoreIdPrefix     = "cache:feedStore:id:"
	cacheFeedStoreUserIdPrefix = "cache:feedStore:userId:"
)

type (
	feedStoreModel interface {
		FeedStoreInsert(ctx context.Context, data *FeedStore, delCacheKeys ...string) (int64, error)
		FeedStoreFindById(ctx context.Context, id int64) (*FeedStore, error)
		FeedStoreUpdateById(ctx context.Context, id int64, updateObj *FeedStore, delCacheKeys []string, fields ...string) (int64, error)
		FeedStoreDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error)
	}

	FeedStore struct {
		Id        int64      `db:"id" gorm:"column:id" json:"id,omitempty"`
		UserId    int64      `db:"user_id" gorm:"column:user_id" json:"user_id,omitempty"`          // 用户ID
		FeedNum   int64      `db:"feed_num" gorm:"column:feed_num" json:"feed_num,omitempty"`       // 私聊数量
		CreatedAt *time.Time `db:"created_at" gorm:"column:created_at" json:"created_at,omitempty"` // 创建时间
		UpdatedAt *time.Time `db:"updated_at" gorm:"column:updated_at" json:"updated_at,omitempty"` // 更新时间
	}
)

func (m *defaultModel) FeedStoreInsert(ctx context.Context, data *FeedStore, delCacheKeys ...string) (int64, error) {
	feedStoreUserIdKey := fmt.Sprintf("%s%v", cacheFeedStoreUserIdPrefix, data.UserId)
	afterDel := true

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))
	delCacheAllKeys = append(delCacheAllKeys, feedStoreUserIdKey)

	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	if data.Id != 0 {
		afterDel = false
		feedStoreIdKey := fmt.Sprintf("%s%v", cacheFeedStoreIdPrefix, data.Id)
		delCacheAllKeys = append(delCacheAllKeys, feedStoreIdKey)
	}

	result, err := m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.WithContext(ctx).Model(&FeedStore{}).Create(data)
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

	if err != nil {
		return 0, err
	}

	if afterDel {
		feedStoreIdKey := fmt.Sprintf("%s%v", cacheFeedStoreIdPrefix, data.Id)
		err = m.db.DelCacheKeysAndDelay(ctx, feedStoreIdKey)
		if err != nil {
			return 0, err
		}
	}
	return result, err

}

func (m *defaultModel) FeedStoreFindById(ctx context.Context, id int64) (*FeedStore, error) {
	feedStoreIdKey := fmt.Sprintf("%s%v", cacheFeedStoreIdPrefix, id)
	var resp FeedStore
	err := m.db.QueryByCtx(ctx, &resp, feedStoreIdKey, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.WithContext(ctx).Model(&FeedStore{}).Where("`id`=?", id).Take(r).Error
	})
	return &resp, err

}

func (m *defaultModel) FeedStoreUpdateById(ctx context.Context, id int64, updateObj *FeedStore, delCacheKeys []string, fields ...string) (int64, error) {
	if updateObj == nil {
		return 0, nil
	}

	data, err := m.FeedStoreFindById(ctx, id)
	if err != nil {
		return 0, err
	}
	feedStoreIdKey := fmt.Sprintf("%s%v", cacheFeedStoreIdPrefix, data.Id)
	feedStoreUserIdKey := fmt.Sprintf("%s%v", cacheFeedStoreUserIdPrefix, data.UserId)

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))

	delCacheAllKeys = append(delCacheAllKeys, feedStoreIdKey)

	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		upTx := db.WithContext(ctx).Model(&FeedStore{}).Where("`id`=?", id)
		if len(fields) > 0 {
			upTx = upTx.Select(strings.Join(fields, ",")).Updates(updateObj)
		} else {
			upTx = upTx.Save(updateObj)
		}
		return upTx.RowsAffected, upTx.Error
	}, delCacheAllKeys...)

}

func (m *defaultModel) FeedStoreDeleteById(ctx context.Context, id int64, delCacheKeys ...string) (int64, error) {

	data, err := m.FeedStoreFindById(ctx, id)
	if err != nil {
		return 0, err
	}

	feedStoreIdKey := fmt.Sprintf("%s%v", cacheFeedStoreIdPrefix, id)
	feedStoreUserIdKey := fmt.Sprintf("%s%v", cacheFeedStoreUserIdPrefix, data.UserId)

	delCacheAllKeys := make([]string, 0, 2+len(delCacheKeys))

	// 0
	delCacheAllKeys = append(delCacheAllKeys, feedStoreIdKey, feedStoreIdKey, feedStoreUserIdKey)

	if len(delCacheKeys) > 0 {
		delCacheAllKeys = append(delCacheAllKeys, delCacheKeys...)
	}

	return m.db.ExecCtx(ctx, func(ctx context.Context, db *gorm.DB) (int64, error) {
		res := db.WithContext(ctx).Where("id = ?", id).Delete(&FeedStore{})
		return res.RowsAffected, res.Error
	}, delCacheAllKeys...)

}

func (FeedStore) TableName() string {
	return "feed_store"
}
