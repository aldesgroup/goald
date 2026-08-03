// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationQuery model
type TranslationQueryModel struct {
	g.IQueryParamsObjectModel
	namespace *g.StringField
	key       *g.StringField
}

// this is the main way to refer to the TranslationQuery model in the applicative code
func TranslationQuery() *TranslationQueryModel {
	return translationQuery
}

// internal variables
var (
	translationQuery     *TranslationQueryModel
	translationQueryOnce sync.Once
)

// fully describing each of this model's properties & relationships
func NewTranslationQueryModel() *TranslationQueryModel {
	thisModel := &TranslationQueryModel{IQueryParamsObjectModel: g.NewQueryParamsObjectModel()}
	thisModel.namespace = g.AddStringField(thisModel, "TranslationQuery", "Namespace", false)
	thisModel.key = g.AddStringField(thisModel, "TranslationQuery", "Key", false)

	return thisModel
}

// making sure the TranslationQuery model exists at app startup
func init() {
	translationQueryOnce.Do(func() {
		translationQuery = NewTranslationQueryModel()
	})

	// this helps dynamically access to the TranslationQuery model
	g.RegisterModel("TranslationQuery", translationQuery)
}

// accessing all the TranslationQuery model's properties and relationships

func (T *TranslationQueryModel) Namespace() *g.StringField {
	return T.namespace
}

func (T *TranslationQueryModel) Key() *g.StringField {
	return T.key
}
