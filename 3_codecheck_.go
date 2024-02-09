// ------------------------------------------------------------------------------------------------
// The code here is about checking that the devs haven't forgotten some stuff, like telling
// each relationship's type, some properties' size, if a class is persisted or not, etc.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"log"
	"time"
)

func (thisServer *server) runCodeChecks() {
	start := time.Now()

	for className, boClass := range getAllClasses() {
		thisServer.checkClass(className, boClass)
	}

	log.Printf("done checking the code in %s", time.Since(start))
}

func (thisServer *server) checkClass(className string, boClass IBusinessObjectClass) {
	nbChildToParentRelationships := 0

	// class-level controls
	if className != ToPascal(className) {
		panicf("The class name '%s' should be pascal-cased, i.e. %s", className, ToPascal(className))
	}

	// checks for the persistency requirements
	if boClass.base().isPersisted() {

		// checking the fields
		for _, field := range boClass.base().fields {
			switch field := field.(type) {
			case *StringField:
				if field.name != "ID" && field.size == 0 {
					panicf("Field '%s.%s' should have a max size set", className, field.name)
				}
			}
		}

		// checking the relationships
		for _, relationship := range boClass.base().relationships {
			if relationship.relationType == 0 {
				panicf("Relationship '%s.%s' should have a defined type, with SetChildToParent(), "+
					"SetSourceToTarget() or SetOneWay()", className, relationship.name)
			}

			if relationship.relationType == relationshipTypeCHILDxTOxPARENT {
				nbChildToParentRelationships++
			}

			if nbChildToParentRelationships > 1 {
				panicf("There cannot be more than one child to parent relationship in '%s'", className)
			}

			if relationship.relationType == relationshipTypePARENTxTOxCHILDREN && !relationship.multiple {
				panicf("We do not handle 1-1 child-parent relationship for now. "+
					"Please re-design relationship '%s.%s'", className, relationship.name)
			}
		}
	}

	// TODO property sizes, when relevant
	// TODO float precision
	// TODO
	// TODO SOON: set primary reference, or none
	// TODO SOON: field / relationshop i/o descriptions
	// TODO SOON: enum & listEnum auto-maxlength
	// TODO
	// TODO	LATER: no column name on not-persisted links
	// TODO LATER: allow custom table name
	// TODO LATER: allow custom column name
	// TODO LATER: tracking policy
	// TODO LATER: unique table name per DB
	// TODO LATER: unique column name per
	// TODO LATER: personal info asserted - with suggestions! (lastname, firstName, mail, email, phone, etc.)
	// TODO LATER: confidential info asserted - with suggestions! (password, pass, passwd)
}
