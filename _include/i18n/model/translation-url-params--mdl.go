// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationUrlParams model
type TranslationUrlParamsModel struct {
	g.IURLQueryParamsModel
	namespace *g.StringField
	key       *g.StringField
}

// this is the main way to refer to the TranslationUrlParams model in the applicative code
func TranslationUrlParams() *TranslationUrlParamsModel {
	return translationUrlParams
}

// internal variables
var (
	translationUrlParams     *TranslationUrlParamsModel
	translationUrlParamsOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewTranslationUrlParamsModel() *TranslationUrlParamsModel {
	thisModel := &TranslationUrlParamsModel{IURLQueryParamsModel: g.NewURLQueryParamsModel()}
	thisModel.namespace = g.AddStringField(thisModel, "TranslationUrlParams", "Namespace", false)
	thisModel.key = g.AddStringField(thisModel, "TranslationUrlParams", "Key", false)

	return thisModel
}

// making sure the TranslationUrlParams model exists at app startup
func init() {
	translationUrlParamsOnce.Do(func() {
		translationUrlParams = NewTranslationUrlParamsModel()
	})

	// this helps dynamically access to the TranslationUrlParams model
	g.RegisterModel("TranslationUrlParams", translationUrlParams)
}

// accessing all the TranslationUrlParams model's properties and relationships

func (T *TranslationUrlParamsModel) Namespace() *g.StringField {
	return T.namespace
}

func (T *TranslationUrlParamsModel) Key() *g.StringField {
	return T.key
}
