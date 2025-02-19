package i18n

import (
	"github.com/aldesgroup/goald"
	class "github.com/aldesgroup/goald/_include/_class"
)

type Translation struct {
	goald.BusinessObject
	Lang      string `json:"-"`
	Namespace string
	Key       string
	Value     string
}

func init() {
	class.Translation().SetNotPersisted()
}
