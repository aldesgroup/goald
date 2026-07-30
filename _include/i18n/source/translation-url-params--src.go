// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TranslationUrlParamsModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForTranslationUrlParams(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &TranslationUrlParamsModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "TranslationUrlParams", lastMod)}
}

func (this *TranslationUrlParamsModelSource) NewObject() any {
	return &i18n.TranslationUrlParams{}
}

func (this *TranslationUrlParamsModelSource) NewSlice() any {
	return []*i18n.TranslationUrlParams{}
}
