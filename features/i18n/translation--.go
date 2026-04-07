package i18n

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/i18n/model"
)

type Translation struct {
	goald.BusinessObject
	Lang      string `json:"-"`
	Namespace string
	Key       string
	Value     string
}

func init() {
	model.Translation().SetNotPersisted()
}
