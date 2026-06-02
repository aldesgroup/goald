package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/i18n/model"
)

type TranslationUrlParams struct {
	g.URLQueryParams
	Namespace string `json:"namespace" io:"o*" desc:"the namespace of the translation keys to retrieve (e.g. 'Common')"`
	Key       string `json:"key"       io:"o*" desc:"the key for the unique translation to retrieve"`
}

func init() {
	model.TranslationUrlParams().SetDescription("The URL query parameters for the endpoint to retrieve translations")
}
