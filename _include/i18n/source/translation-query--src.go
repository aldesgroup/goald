// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TranslationQueryModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForTranslationQuery(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &TranslationQueryModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "TranslationQuery", lastMod)}
}

func (this *TranslationQueryModelSource) NewObject() any {
	return &i18n.TranslationQuery{}
}

func (this *TranslationQueryModelSource) NewSlice() any {
	return []*i18n.TranslationQuery{}
}
