// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationUrlParams class
type translationUrlParamsClass struct {
	g.IURLQueryParamsClass
	namespace *g.StringField
	key       *g.StringField
}

// this is the main way to refer to the TranslationUrlParams class in the applicative code
func TranslationUrlParams() *translationUrlParamsClass {
	return translationUrlParams
}

// internal variables
var (
	translationUrlParams     *translationUrlParamsClass
	translationUrlParamsOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationUrlParamsClass() *translationUrlParamsClass {
	newClass := &translationUrlParamsClass{IURLQueryParamsClass: g.NewURLQueryParamsClass()}
	newClass.namespace = g.NewStringField(newClass, "Namespace", false)
	newClass.key = g.NewStringField(newClass, "Key", false)

	return newClass
}

// making sure the TranslationUrlParams class exists at app startup
func init() {
	translationUrlParamsOnce.Do(func() {
		translationUrlParams = newTranslationUrlParamsClass()
	})

	// this helps dynamically access to the TranslationUrlParams class
	g.RegisterClass("TranslationUrlParams", translationUrlParams)
}

// accessing all the TranslationUrlParams class' properties and relationships

func (t *translationUrlParamsClass) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationUrlParamsClass) Key() *g.StringField {
	return t.key
}
