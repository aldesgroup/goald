// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationUrlParams specs
type translationUrlParamsSpecs struct {
	g.IURLQueryParamsSpecs
	namespace *g.StringField
	key       *g.StringField
}

// this is the main way to refer to the TranslationUrlParams specs in the applicative code
func TranslationUrlParams() *translationUrlParamsSpecs {
	return translationUrlParams
}

// internal variables
var (
	translationUrlParams     *translationUrlParamsSpecs
	translationUrlParamsOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationUrlParamsSpecs() *translationUrlParamsSpecs {
	newSpecs := &translationUrlParamsSpecs{IURLQueryParamsSpecs: g.NewURLQueryParamsSpecs()}
	newSpecs.namespace = g.NewStringField(newSpecs, "Namespace", false)
	newSpecs.key = g.NewStringField(newSpecs, "Key", false)

	return newSpecs
}

// making sure the TranslationUrlParams specs exists at app startup
func init() {
	translationUrlParamsOnce.Do(func() {
		translationUrlParams = newTranslationUrlParamsSpecs()
	})

	// this helps dynamically access to the TranslationUrlParams specs
	g.RegisterSpecs("TranslationUrlParams", translationUrlParams)
}

// accessing all the TranslationUrlParams class' properties and relationships

func (t *translationUrlParamsSpecs) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationUrlParamsSpecs) Key() *g.StringField {
	return t.key
}
