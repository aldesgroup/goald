// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TranslationUrlParamsClassUtils struct {
	goald.IClassUtilsCore
}

func ClassUtilsForTranslationUrlParams(srcPath, lastMod string) goald.IClassUtils {
	return &TranslationUrlParamsClassUtils{IClassUtilsCore: goald.NewClassUtilsCore(srcPath, lastMod)}
}

func (thisTranslationUrlParamsClassUtils *TranslationUrlParamsClassUtils) NewObject() any {
	return &i18n.TranslationUrlParams{}
}

func (thisTranslationUrlParamsClassUtils *TranslationUrlParamsClassUtils) NewSlice() any {
	return []*i18n.TranslationUrlParams{}
}
