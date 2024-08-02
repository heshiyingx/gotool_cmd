package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

type (
	// FeedStoreDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customFeedStoreModel.
	FeedStoreDBInterface interface {
		feedStoreModel
	}
)
