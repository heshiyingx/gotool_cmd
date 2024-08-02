package code

import (
	"github.com/heshiyingx/gotool/dbext/gormdb"
)

type (
	// ChickenDBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in customChickenModel.
	ChickenDBInterface interface {
		chickenModel
	}
)
