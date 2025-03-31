// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Translation specs
type translationSpecs struct {
	g.IBusinessObjectSpecs
	lang      *g.StringField
	namespace *g.StringField
	key       *g.StringField
	value     *g.StringField
}

// this is the main way to refer to the Translation specs in the applicative code
func Translation() *translationSpecs {
	return translation
}

// internal variables
var (
	translation     *translationSpecs
	translationOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationSpecs() *translationSpecs {
	newSpecs := &translationSpecs{IBusinessObjectSpecs: g.NewBusinessObjectSpecs()}
	newSpecs.lang = g.NewStringField(newSpecs, "Lang", false)
	newSpecs.namespace = g.NewStringField(newSpecs, "Namespace", false)
	newSpecs.key = g.NewStringField(newSpecs, "Key", false)
	newSpecs.value = g.NewStringField(newSpecs, "Value", false)

	return newSpecs
}

// making sure the Translation specs exists at app startup
func init() {
	translationOnce.Do(func() {
		translation = newTranslationSpecs()
	})

	// this helps dynamically access to the Translation specs
	g.RegisterSpecs("Translation", translation)
}

// accessing all the Translation class' properties and relationships

func (t *translationSpecs) Lang() *g.StringField {
	return t.lang
}

func (t *translationSpecs) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationSpecs) Key() *g.StringField {
	return t.key
}

func (t *translationSpecs) Value() *g.StringField {
	return t.value
}
