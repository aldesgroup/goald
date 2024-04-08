package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the TranslationKey class
type translationKeyClass struct {
	g.IBusinessObjectClass
	route        *g.StringField
	part         *g.StringField
	key          *g.StringField
	translations *g.Relationship
}

// this is the main way to refer to the TranslationKey class in the applicative code
func TranslationKey() *translationKeyClass {
	return translationKey
}

// internal variables
var (
	translationKey     *translationKeyClass
	translationKeyOnce sync.Once
)

// fully describing each of this class' properties & relationships
func newTranslationKeyClass() *translationKeyClass {
	newClass := &translationKeyClass{IBusinessObjectClass: g.NewClass()}
	newClass.route = g.NewStringField(newClass, "Route", false)
	newClass.part = g.NewStringField(newClass, "Part", false)
	newClass.key = g.NewStringField(newClass, "Key", false)
	newClass.translations = g.NewRelationship(newClass, "Translations", true, translation)

	return newClass
}

// making sure the TranslationKey class exists at app startup
func init() {
	translationKeyOnce.Do(func() {
		translationKey = newTranslationKeyClass()
	})

	// this helps dynamically access to the TranslationKey class
	g.RegisterClass("TranslationKey", translationKey)
}

// accessing all the TranslationKey class' properties and relationships

func (t *translationKeyClass) Route() *g.StringField {
	return t.route
}

func (t *translationKeyClass) Part() *g.StringField {
	return t.part
}

func (t *translationKeyClass) Key() *g.StringField {
	return t.key
}

func (t *translationKeyClass) Translations() *g.Relationship {
	return t.translations
}
