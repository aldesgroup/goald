// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
"github.com/aldesgroup/goald/features/i18n"
)

type TranslationUrlParamsClass struct {
	goald.IClassCore
}

func ClassForTranslationUrlParams(srcPath, lastMod string) goald.IClass {
	return &TranslationUrlParamsClass{IClassCore: goald.NewClassCore(srcPath, "TranslationUrlParams", lastMod)}
}

func (thisClass *TranslationUrlParamsClass) NewObject() any {
	return &i18n.TranslationUrlParams{}
}

func (thisClass *TranslationUrlParamsClass) NewSlice() any {
	return []*i18n.TranslationUrlParams{}
}
