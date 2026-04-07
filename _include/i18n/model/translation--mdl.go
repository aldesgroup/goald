// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Translation model
type translationModel struct {
	g.IBusinessObjectModel
	lang      *g.StringField
	namespace *g.StringField
	key       *g.StringField
	value     *g.StringField
}

// this is the main way to refer to the Translation model in the applicative code
func Translation() *translationModel {
	return translation
}

// internal variables
var (
	translation     *translationModel
	translationOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationModel() *translationModel {
	newModel := &translationModel{IBusinessObjectModel: g.NewBusinessObjectModel()}
	newModel.lang = g.NewStringField(newModel, "Lang", false)
	newModel.namespace = g.NewStringField(newModel, "Namespace", false)
	newModel.key = g.NewStringField(newModel, "Key", false)
	newModel.value = g.NewStringField(newModel, "Value", false)

	return newModel
}

// making sure the Translation model exists at app startup
func init() {
	translationOnce.Do(func() {
		translation = newTranslationModel()
	})

	// this helps dynamically access to the Translation model
	g.RegisterModel("Translation", translation)
}

// accessing all the Translation class' properties and relationships

func (t *translationModel) Lang() *g.StringField {
	return t.lang
}

func (t *translationModel) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationModel) Key() *g.StringField {
	return t.key
}

func (t *translationModel) Value() *g.StringField {
	return t.value
}
