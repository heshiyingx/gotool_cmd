package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

type (
	// FeedStoreHistoryDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customFeedStoreHistoryModel.
	FeedStoreHistoryDBInterface interface {
		feedStoreHistoryModel
	}
)
