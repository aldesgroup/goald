// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Translation model
type TranslationModel struct {
	g.IBusinessObjectModel
	lang      *g.StringField
	namespace *g.StringField
	key       *g.StringField
	value     *g.StringField
}

// this is the main way to refer to the Translation model in the applicative code
func Translation() *TranslationModel {
	return translation
}

// internal variables
var (
	translation     *TranslationModel
	translationOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewTranslationModel() *TranslationModel {
	thisModel := &TranslationModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	thisModel.lang = g.AddStringField(thisModel, "Translation", "Lang", false)
	thisModel.namespace = g.AddStringField(thisModel, "Translation", "Namespace", false)
	thisModel.key = g.AddStringField(thisModel, "Translation", "Key", false)
	thisModel.value = g.AddStringField(thisModel, "Translation", "Value", false)

	return thisModel
}

// making sure the Translation model exists at app startup
func init() {
	translationOnce.Do(func() {
		translation = NewTranslationModel()
	})

	// this helps dynamically access to the Translation model
	g.RegisterModel("Translation", translation)
}

// accessing all the Translation model's properties and relationships

func (T *TranslationModel) Lang() *g.StringField {
	return T.lang
}

func (T *TranslationModel) Namespace() *g.StringField {
	return T.namespace
}

func (T *TranslationModel) Key() *g.StringField {
	return T.key
}

func (T *TranslationModel) Value() *g.StringField {
	return T.value
}
