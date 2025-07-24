// ------------------------------------------------------------------------------------------------
// The code here is about checking that the devs haven't forgotten some stuff, like telling
// each relationship's type, some properties' size, if a class is persisted or not, etc.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"time"

	core "github.com/aldesgroup/corego"
)

func (thisServer *server) runCodeChecks() {
	start := time.Now()

	for clsName, boSpecs := range specsRegistry.items {
		thisServer.checkSpecs(clsName, boSpecs)
	}

	slog.Info(fmt.Sprintf("done checking the code in %s", time.Since(start)))
}

func (thisServer *server) checkSpecs(clsName className, boSpecs IBusinessObjectSpecs) {
	nbChildToParentRelationships := 0

	// class-level controls
	if expected := core.ToPascal(string(clsName)); string(clsName) != expected {
		core.PanicMsg("The class name '%s' should be pascal-cased, i.e. %s", clsName, expected)
	}

	if !boSpecs.base().abstract {
		// various check, whether there's persistence or not
		for _, field := range boSpecs.base().fields {
			if enumField, ok := field.(*EnumField); ok {
				for _, restrictedValue := range enumField.onlyValues {
					if fmt.Sprintf("%T", restrictedValue) != enumField.enumName {
						core.PanicMsg("Cannot use '%v' (%T) as a '%s' value in class '%s'!",
							restrictedValue, restrictedValue, enumField.enumName, clsName)
					}
				}
			}
		}

		// checks for the persistency requirements
		if boSpecs.base().isPersisted() {
			// checking there's an actual DB configured for this BO class
			if boSpecs.getInDB() == nil {
				core.PanicMsg("Class '%s' should be SetNotPersisted, SetAbstract, or associated with a DB", clsName)
			}

			// checking the fields
			for _, field := range boSpecs.base().fields {
				switch field := field.(type) {
				case *StringField:
					if field.name != "ID" && field.size == 0 && !field.isNotPersisted() {
						core.PanicMsg("Field '%s.%s' should have a max size set, or be SetNotPersisted()", clsName, field.name)
					}
				}
			}
		}

		// checking the relationships
		if boSpecs.base().isPersisted() || boSpecs.base().usedInNativeApp || boSpecs.base().usedInWebApp {
			for _, relationship := range boSpecs.base().relationships {
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
