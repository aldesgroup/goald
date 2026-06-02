// ------------------------------------------------------------------------------------------------
// The code here is about checking that the devs haven't forgotten some stuff, like telling
// each relationship's type, some properties' size, if a class is persisted or not, etc.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
)

// ------------------------------------------------------------------------------------------------
// This function is called in codegen mode, and checks the code of the BO models and endpoints
// ------------------------------------------------------------------------------------------------

func (thisServer *server) runCodeChecks() {
	start := time.Now()

	// checking each BO model
	for clsName, model := range modelRegistry.items {
		thisServer.checkModel(clsName, model)
	}

	// checking each endpoint
	for _, ep := range restRegistry.endpoints {
		thisServer.checkEndpoint(ep)
	}

	slog.Info(fmt.Sprintf("done checking the code in %s", time.Since(start)))
}

// ------------------------------------------------------------------------------------------------
// BO model checking
// ------------------------------------------------------------------------------------------------

func (thisServer *server) checkModel(clsName className, model IBusinessObjectModel) {
	nbChildToParentRelationships := 0

	// class-level controls
	if expected := core.ToPascal(string(clsName)); string(clsName) != expected {
		core.PanicMsg("The model name '%s' should be pascal-cased, i.e. %s", clsName, expected)
	}

	if !model.base().abstract && model.base().description == "" {
		core.PanicMsg("Model '%s' should have a description", clsName)
	}

	if !model.base().abstract {
		// various check, whether there's persistence or not
		for _, field := range model.base().fields {

			// enum-related checks
			if enumField, ok := field.(*EnumField); ok {
				for _, restrictedValue := range enumField.onlyValues {
					if fmt.Sprintf("%T", restrictedValue) != enumField.enumName {
						core.PanicMsg("Cannot use '%v' (%T) as a '%s' value in model '%s'!",
							restrictedValue, restrictedValue, enumField.enumName, clsName)
					}
				}
			}

			// generic checks - JSON name
			thisServer.genericPropertyCodeCheck(field)
		}

		// checks for the persistency requirements
		if model.base().isPersisted() {
			// checking there's an actual DB configured for this BO model
			if model.getInDB() == nil {
				core.PanicMsg("Model '%s' should be SetNotPersisted, SetAbstract, or associated with a DB", clsName)
			}

			// checking the fields
			for _, field := range model.base().fields {
				// type-related checks
				switch field := field.(type) {
				case *StringField:
					if field.name != "ID" && field.size == 0 && !field.isNotPersisted() {
						core.PanicMsg("Field '%s.%s' should have a max size set, or be SetNotPersisted()", clsName, field.name)
					}
				}
			}
		}

		// checking the relationships - generic checks
		for _, relationship := range model.base().relationships {
			thisServer.genericPropertyCodeCheck(relationship)
		}

		// checking the relationships - when there's proven I/O with an app or databases
		if model.base().isPersisted() || model.base().usedInNativeApp || model.base().usedInWebApp {
			for _, relationship := range model.base().relationships {
				if relationship.relationType == 0 {
					core.PanicMsg("Relationship '%s.%s' should have a defined type, with SetChildToParent(), "+
						"SetSourceToTarget() or SetOneWay()", clsName, relationship.name)
				}

				if relationship.relationType == relationshipTypeCHILDxTOxPARENT {
					nbChildToParentRelationships++
				}

				if nbChildToParentRelationships > 1 {
					core.PanicMsg("There cannot be more than one child to parent relationship in '%s'", clsName)
				}
			}
		}
	}

	// TODO property sizes, when relevant
	// TODO float precision
	// TODO
	// TODO SOON: set primary reference, or none
	// TODO SOON: field / relationshop i/o descriptions
	// TODO SOON: enum & listEnum auto-maxlength
	// TODO SOON: query BObj : prevent some property types
	// TODO
	// TODO	LATER: no column name on not-persisted links
	// TODO LATER: allow custom table name
	// TODO LATER: allow custom column name
	// TODO LATER: tracking policy
	// TODO LATER: unique table name per DB
	// TODO LATER: unique column name per property
	// TODO LATER: personal info asserted - with suggestions! (lastname, firstName, mail, email, phone, etc.)
	// TODO LATER: confidential info asserted - with suggestions! (password, pass, passwd)
}

var ioTagsMap = map[string]string{
	"in": "non-mandatory input property",
	"i*": "mandatory input property",
	"o*": "pure computed / output property",
}

var ioTagsStr = core.MapToString(ioTagsMap, true, ": ", ",\n")

func (thisServer *server) genericPropertyCodeCheck(property iBusinessObjectProperty) {
	// Empty ("-") or valid camel-case JSON name
	jsonTags := strings.Split(property.getTag("json"), ",")
	jsonName := jsonTags[0]
	if jsonName == "" || jsonName != "-" && jsonName != core.PascalToCamel(property.getName()) {
		core.PanicMsg("Property '%s.%s' should be ignored with \"-\", or have a json tag set to '%s', not '%s'",
			property.ownerModel().base().name, property.getName(), core.PascalToCamel(property.getName()), jsonName)
	}

	// Valid I/O tag
	ioTag := property.getTag("io")
	if _, ok := ioTagsMap[ioTag]; !ok {
		core.PanicMsg("Property '%s.%s' should have an 'io' tag equals to '%s' but should have one of these values: \n%s",
			property.ownerModel().base().name, property.getName(), ioTag, ioTagsStr)
	}

	// Non-empty description in the "desc" tag
	desc := property.getTag("desc")
	if desc == "" {
		core.PanicMsg("Property '%s.%s' should have a non-empty description in the 'desc' tag", property.ownerModel().base().name, property.getName())
	}
}

// ------------------------------------------------------------------------------------------------
// Endpoint checking
// ------------------------------------------------------------------------------------------------

func (thisServer *server) checkEndpoint(endpoint iEndpoint) {
	if endpoint.getLabel() == "" {
		core.PanicMsg("Endpoint '%s' should have a non-empty label", endpoint.getPathAsString())
	}
	if endpoint.getDescription() == "" {
		core.PanicMsg("Endpoint '%s' should have a non-empty description", endpoint.getPathAsString())
	}
}
