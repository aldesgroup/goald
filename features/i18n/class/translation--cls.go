// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
"github.com/aldesgroup/goald/features/i18n"
)

type TranslationClass struct {
	goald.IClassCore
}

func ClassForTranslation(srcPath, lastMod string) goald.IClass {
	return &TranslationClass{IClassCore: goald.NewClassCore(srcPath, "Translation", lastMod)}
}

func (thisClass *TranslationClass) NewObject() any {
	return &i18n.Translation{}
}

func (thisClass *TranslationClass) NewSlice() any {
	return []*i18n.Translation{}
}
