package code

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type (
	// FeedStoreDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customFeedStoreModel.
	FeedStoreDBInterface interface {
		feedStoreModel
		FeedStoreNoCacheFindOneByUserId(ctx context.Context, userId int64) (*FeedStore, error)
		FeedStoreNoCacheNolimitFindOneByUserId(ctx context.Context, userId int64) (*FeedStore, error)
	}
)

func (m *defaultModel) FeedStoreNoCacheFindOneByUserId(ctx context.Context, userId int64) (*FeedStore, error) {
	feedStoreUserIdKey := fmt.Sprintf("%s%v", cacheFeedStoreUserIdPrefix, userId)

	var resp FeedStore
	err := m.db.QuerySingleNoCacheCtx(ctx, feedStoreUserIdKey, &resp, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.Model(&FeedStore{}).Where("`user_id`= ?", userId).Take(r).Error
	})
	return &resp, err

}
func (m *defaultModel) FeedStoreNoCacheNolimitFindOneByUserId(ctx context.Context, userId int64) (*FeedStore, error) {
	//feedStoreUserIdKey := fmt.Sprintf("%s%v", cacheFeedStoreUserIdPrefix, userId)

	var resp FeedStore
	err := m.db.QueryNoCacheCtx(ctx, &resp, func(ctx context.Context, r any, db *gorm.DB) error {
		return db.Model(&FeedStore{}).Where("`user_id`= ?", userId).Take(r).Error
	})
	return &resp, err

}