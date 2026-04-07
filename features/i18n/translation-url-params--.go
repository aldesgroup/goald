package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/i18n/model"
)

type TranslationUrlParams struct {
	g.URLQueryParams
	Namespace string
	Key       string
}

func init() {
	model.TranslationUrlParams().Namespace().SetMandatory()
}
