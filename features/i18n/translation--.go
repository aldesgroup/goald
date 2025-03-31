package i18n

import (
	"github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type Translation struct {
	goald.BusinessObject
	Lang      string `json:"-"`
	Namespace string
	Key       string
	Value     string
}

func init() {
	specs.Translation().SetNotPersisted()
}
