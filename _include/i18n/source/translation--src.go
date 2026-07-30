// Generated file, do not edit!
package source

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

type TranslationModelSource struct {
	goald.IBusinessObjectModelSource
}

func ForTranslation(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &TranslationModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "Translation", lastMod)}
}

func (this *TranslationModelSource) NewObject() any {
	return &i18n.Translation{}
}

func (this *TranslationModelSource) NewSlice() any {
	return []*i18n.Translation{}
}
