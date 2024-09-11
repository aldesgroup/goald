package i18n

import (
	g "github.com/aldesgroup/goald"
	class "github.com/aldesgroup/goald/_include/_class"
)

type TranslationUrlParams struct {
	g.URLQueryParams
	Namespace string
	Key       string
}

func init() {
	class.TranslationUrlParams().Namespace().SetMandatory()
}
