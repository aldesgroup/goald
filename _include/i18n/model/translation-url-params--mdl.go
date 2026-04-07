// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationUrlParams model
type translationUrlParamsModel struct {
	g.IURLQueryParamsModel
	namespace *g.StringField
	key       *g.StringField
}

// this is the main way to refer to the TranslationUrlParams model in the applicative code
func TranslationUrlParams() *translationUrlParamsModel {
	return translationUrlParams
}

// internal variables
var (
	translationUrlParams     *translationUrlParamsModel
	translationUrlParamsOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationUrlParamsModel() *translationUrlParamsModel {
	newModel := &translationUrlParamsModel{IURLQueryParamsModel: g.NewURLQueryParamsModel()}
	newModel.namespace = g.NewStringField(newModel, "Namespace", false)
	newModel.key = g.NewStringField(newModel, "Key", false)

	return newModel
}

// making sure the TranslationUrlParams model exists at app startup
func init() {
	translationUrlParamsOnce.Do(func() {
		translationUrlParams = newTranslationUrlParamsModel()
	})

	// this helps dynamically access to the TranslationUrlParams model
	g.RegisterModel("TranslationUrlParams", translationUrlParams)
}

// accessing all the TranslationUrlParams class' properties and relationships

func (t *translationUrlParamsModel) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationUrlParamsModel) Key() *g.StringField {
	return t.key
}
