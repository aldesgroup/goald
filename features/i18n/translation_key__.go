package i18n

import "github.com/aldesgroup/goald"

type TranslationKey struct {
	goald.BusinessObject
	Route        string
	Part         string
	Key          string
	Translations []*Translation
}
