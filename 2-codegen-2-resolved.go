// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const resolvedTEMPLATE = `// Generated file, do not edit!
package include

import (
	%s
)

func init() {
	// registrering all the implementations for the interface-type Business Objects that other BOs link to through relationships
	%s
}

`

const resolvedFILE = "_include/resolved.go"

func (thisServer *server) generateResolvedRelationships(srcdir string, regen bool) (codeChanged bool) {
	// we need to do it only if we are regenerating, or if the file does not exist yet
	if regen || !core.FileExists(resolvedFILE) {

		imports := map[string]string{}
		registrations := []string{}
		implementations := map[string][]IClass{}

		// looking for the polymorphic relationships that we have
		for _, model := range modelRegistry.items {
			for _, relationship := range model.base().relationships {
				if relationship.polymorphic {
					sourceObjCls := getClass(model)                                         // e.g. ClassForSourceObj
					sourceObject := sourceObjCls.NewObject()                                // e.g.: *SourceObj
					sourceObjTyp := utils.TypeOf(sourceObject, true)                        // e.g. Type SourceObj
					targetFldTyp := sourceObjTyp.FieldByName(relationship.getName()).Type() // e.g. Type ITargetObj or []ITargetObj

					if relationship.multiple {
						targetFldTyp = targetFldTyp.Elem() // e.g. Type ITargetObj
					}

					// finding all the classes implementing the interface-type BO
					if _, alreadyHandled := implementations[targetFldTyp.Name()]; !alreadyHandled {
						implementations[targetFldTyp.Name()] = findImplementionsOfInterface(targetFldTyp)
					}

					// handling the model for each class
					models := []string{}
					for _, concreteTargetObjCls := range implementations[targetFldTyp.Name()] {
						if imports[concreteTargetObjCls.getPackage()] == "" {
							imports[concreteTargetObjCls.getPackage()] = getImportLineForClass(concreteTargetObjCls) // the needed import
						}
						models = append(models, fmt.Sprintf("%s_model.%s()", concreteTargetObjCls.getPackage(), concreteTargetObjCls.getClassName())) // e.g. packagename_model.ConcreteTargetObj
					}

					// adding the "resolution" of the interface type with the corresponding implementations
					if len(models) > 0 {
						registrations = append(registrations, fmt.Sprintf("%s_model.%s().%s().SetTargets(%s)",
							sourceObjCls.getPackage(), sourceObjCls.getClassName(), relationship.name, strings.Join(models, ", ")))

						// let's make sure the source class, i.e. the class owning the relationship, can be imported
						if imports[sourceObjCls.getPackage()] == "" {
							imports[sourceObjCls.getPackage()] = getImportLineForClass(sourceObjCls) // the needed import
						}
					}
				}
			}
		}

		if len(registrations) > 0 {
			core.WriteStringToFile(path.Join(srcdir, resolvedFILE), resolvedTEMPLATE,
				strings.Join(core.GetSortedValues(imports), "\n\t"), strings.Join(registrations, "\n\t"))
		}
	}

	return
}

// this function finds all the implementations of a given interface
func findImplementionsOfInterface(interfaceType utils.GoaldType) (implementations []IClass) {
	// browsing through all the non-interface classes to find the implementations
	for _, class := range classRegistry.items {
		if !class.isInterface() {
			if boType := utils.TypeOf(class.NewObject(), false); boType.Implements(interfaceType) {
				implementations = append(implementations, class)
			}
		}
	}

	return
}
