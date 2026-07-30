// Generated file, do not edit!
package i18n

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_include/i18n/model"
)

type Translation struct {
	goald.BusinessObject
	Lang      string `json:"lang,omitempty"      io:"i*" desc:"the language code (e.g. 'en', 'fr', 'de', etc.)"`
	Namespace string `json:"namespace,omitempty" io:"in" desc:"the namespace of the translation"`
	Key       string `json:"key,omitempty"       io:"o*" desc:"the key of the translation"`
	Value     string `json:"value,omitempty"     io:"o*" desc:"the value of the translation"`
}

func init() {
	model.Translation().SetDescription("A text in a given language, identified by a namespace and a key")
	model.Translation().SetNotPersisted()
}
