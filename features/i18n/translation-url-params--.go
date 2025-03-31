package i18n

import (
	g "github.com/aldesgroup/goald"
	specs "github.com/aldesgroup/goald/_include/_specs"
)

type TranslationUrlParams struct {
	g.URLQueryParams
	Namespace string
	Key       string
}

func init() {
	specs.TranslationUrlParams().Namespace().SetMandatory()
}
