package i18n

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_generated/class"
)

type Translation struct {
	goald.BusinessObject
	Lang  string `json:"-"`
	Route string
	Part  string
	Key   string
	Value string
}

func init() { //
	class.Translation().SetNotPersisted()
}
