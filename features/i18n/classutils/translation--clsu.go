// Generated file, do not edit!
package classutils

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TranslationClassUtils struct {
	goald.IClassUtilsCore
}

func ClassUtilsForTranslation(srcPath, lastMod string) goald.IClassUtils {
	return &TranslationClassUtils{IClassUtilsCore: goald.NewClassUtilsCore(srcPath, lastMod)}
}

func (thisUtils *TranslationClassUtils) NewObject() any {
	return &i18n.Translation{}
}

func (thisUtils *TranslationClassUtils) NewSlice() any {
	return []*i18n.Translation{}
}
