package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_generated/class"
)

type TranslationUrlParams struct {
	g.URLQueryParams
	Route string
	Part  string
	Key   string
}

func init() {
	class.TranslationUrlParams().Route().SetMandatory()
}
