package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

type (
	// NutritionStoreHistoryDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customNutritionStoreHistoryModel.
	NutritionStoreHistoryDBInterface interface {
		nutritionStoreHistoryModel
	}
)
