package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

type (
	// NutritionStoreDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customNutritionStoreModel.
	NutritionStoreDBInterface interface {
		nutritionStoreModel
	}
)
