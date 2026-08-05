// Generated file, do not edit!
package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/i18n/model"
)

type TranslationQuery struct {
	g.SearchParamValues
	Namespace string `json:"namespace,omitempty" io:"o*" desc:"the namespace of the translation keys to retrieve (e.g. 'Common')"`
	Key       string `json:"key,omitempty"       io:"o*" desc:"the key for the unique translation to retrieve"`
}

func init() {
	model.TranslationQuery().SetDescription("The URL query parameters for the endpoint to retrieve translations")
}
