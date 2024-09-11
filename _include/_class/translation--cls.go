// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the Translation class
type translationClass struct {
	g.IBusinessObjectClass
	lang      *g.StringField
	namespace *g.StringField
	key       *g.StringField
	value     *g.StringField
}

// this is the main way to refer to the Translation class in the applicative code
func Translation() *translationClass {
	return translation
}

// internal variables
var (
	translation     *translationClass
	translationOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationClass() *translationClass {
	newClass := &translationClass{IBusinessObjectClass: g.NewClass()}
	newClass.lang = g.NewStringField(newClass, "Lang", false)
	newClass.namespace = g.NewStringField(newClass, "Namespace", false)
	newClass.key = g.NewStringField(newClass, "Key", false)
	newClass.value = g.NewStringField(newClass, "Value", false)

	return newClass
}

// making sure the Translation class exists at app startup
func init() {
	translationOnce.Do(func() {
		translation = newTranslationClass()
	})

	// this helps dynamically access to the Translation class
	g.RegisterClass("Translation", translation)
}

// accessing all the Translation class' properties and relationships

func (t *translationClass) Lang() *g.StringField {
	return t.lang
}

func (t *translationClass) Namespace() *g.StringField {
	return t.namespace
}

func (t *translationClass) Key() *g.StringField {
	return t.key
}

func (t *translationClass) Value() *g.StringField {
	return t.value
}
