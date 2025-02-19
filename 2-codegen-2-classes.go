// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	u "github.com/aldesgroup/goald/features/utils"
)

const classTEMPLATE = `// Generated file, do not edit!
package class

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the $$Upper$$ class
type $$lower$$Class struct {
$$propdecl$$
}

// this is the main way to refer to the $$Upper$$ class in the applicative code
func $$Upper$$() *$$lower$$Class {
	return $$lower$$
}

// internal variables
var $$lower$$ *$$lower$$Class
var $$lower$$Once sync.Once

// fully describing each of this class' properties & relationships
func new$$Upper$$Class() *$$lower$$Class {
	$$propinit$$

	return newClass
}

// making sure the $$Upper$$ class exists at app startup
func init() {
	$$lower$$Once.Do(func() {
		$$lower$$ = new$$Upper$$Class()
	})

	// this helps dynamically access to the $$Upper$$ class
	g.RegisterClass("$$Upper$$", $$lower$$)
}

// accessing all the $$Upper$$ class' properties and relationships

$$accessors$$

`

const classFOLDER = "_include/_class"
const classFILExSUFFIX = "--cls.go"
const classFILExSUFFIXxLEN = len(classFILExSUFFIX)
const classNAMExSUFFIX = "Class"
const newline = "\n"

func (thisServer *server) generateObjectClasses(srcdir string, regen bool) (codeChanged bool) {
	// a type just used here
	type classFile struct {
		modTime  time.Time
		filename string
	}

	// where the class files will be generated
	classDir := u.EnsureDir(srcdir, classFOLDER)

	// we'll gather all the existing class files
	existingClassFiles := map[className]*classFile{}

	// so, let's read the class folder
	classEntries, errDir := os.ReadDir(classDir)
	u.PanicErrf(errDir, "Could not read the class folder")
	for _, classEntry := range classEntries {
		classEntryInfo, errInfo := classEntry.Info()
		u.PanicErrf(errInfo, "Could not read info for file '%s'", classEntry.Name())
		classEntryName := className(u.KebabToPascal(classEntry.Name()[:len(classEntry.Name())-classFILExSUFFIXxLEN]))
		existingClassFiles[classEntryName] = &classFile{
			modTime:  classEntryInfo.ModTime(),
			filename: classEntry.Name(),
		}
	}

	// let's see what we have in terms of business objects
	for name, classUtils := range classUtilsRegistry.content {
		// considering only the business objects of THIS module
		if classUtils.getModule() == getCurrentModuleName() {
			// do we need to regen the class file?
			if existingClass := existingClassFiles[name]; regen ||
				existingClass == nil || existingClass.modTime.Before(classUtils.getLastBOMod()) {
				// generating the missing or outdated class
				generateObjectClass(classDir, classUtils)

				// the code has changed
				codeChanged = true
			}

			// flagging this business object class as NOT unneeded (i.e. needed)
			delete(existingClassFiles, name)
		}
	}

	// removing the unneeded classes
	for _, unneededClass := range existingClassFiles {
		slog.Info(fmt.Sprintf("removing %s", unneededClass.filename))
		if errRem := os.Remove(path.Join(classDir, unneededClass.filename)); errRem != nil {
			u.PanicErrf(errRem, "Could not delete class file '%s'", unneededClass.filename)
		}
	}

	return
}

type classGenContext struct {
	superType     *u.GoaldType
	propertyNames []string
	propertiesMap map[string]*classGenPropertyInfo
}

type classGenPropertyInfo struct {
	propType   u.TypeFamily
	multiple   bool
	targetType string
}

func generateObjectClass(classDir string, classUtils IClassUtils) {
	// starting to build the file content, with the same context
	context := &classGenContext{propertiesMap: map[string]*classGenPropertyInfo{}}

	// trivial filling of the template
	class := string(classUtils.getClass())
	content := strings.ReplaceAll(classTEMPLATE, "$$Upper$$", class)
	content = strings.ReplaceAll(content, "$$lower$$", u.PascalToCamel(class))

	// declaring the properties of the classe
	content = strings.Replace(content, "$$propdecl$$", buildPropDecl(classUtils, context), 1)

	// valueing the properties
	content = strings.Replace(content, "$$propinit$$", buildPropInit(classUtils, context), 1)

	// building the accessors to the properties
	content = strings.Replace(content, "$$accessors$$", buildAccessors(classUtils, context), 1)

	// writing to file
	u.WriteToFile(content, classDir, u.PascalToKebab(class)+classFILExSUFFIX)

	slog.Info(fmt.Sprintf("(Re-)generated class %s", class))
}

func buildPropDecl(classUtils IClassUtils, context *classGenContext) (result string) {
	// getting the object's type
	bObjType := u.TypeOf(classUtils.NewObject(), true)

	// the very first property, field #0, MUST be the business object's super class
	superClassField := bObjType.Field(0)
	if !superClassField.IsAnonymous() || !u.PointerTo(superClassField.Type()).Implements(typeIxBUSINESSxOBJECT) {
		u.Panicf("%s: this object's first property should be the BO it inherits from, i.e."+
			"goald.BusinessObject, or one of its descendants", classUtils.getClass())
	}

	if context.superType = superClassField.Type(); context.superType.Equals(typeBUSINESSxOBJECT) {
		result += "g.IBusinessObjectClass"
	} else if context.superType.Equals(typeURLxQUERYxOBJECT) {
		result += "g.IURLQueryParamsClass"
	} else {
		result += "" + u.PascalToCamel(superClassField.Type().Name()) + classNAMExSUFFIX
	}

	// browsing the entity's properties
	for fieldNum := 1; fieldNum < bObjType.NumField(); fieldNum++ {
		// getting the current field
		field := bObjType.Field(fieldNum)

		// detecting its type and multiplicity
		typeFamily, multiple := u.GetTypeFamily(field, typeIxBUSINESSxOBJECT, typeIxENUM)

		// adding to the context, and the class file content
		if typeFamily != u.TypeFamilyUNKNOWN {
			context.propertyNames = append(context.propertyNames, field.Name()) // we're keeping the original order

			targetType := ""                            // makes no sense for BO fields...
			if typeFamily == u.TypeFamilyRELATIONSHIP { // ... but it does for relationships
				entityType := u.IfThenElse(multiple, field.Type().Elem(), field.Type())
				targetType = entityType.Elem().Name()
			} else if typeFamily == u.TypeFamilyENUM { // or enums.
				targetType = field.Type().String()
			}

			context.propertiesMap[field.Name()] = &classGenPropertyInfo{typeFamily, multiple, targetType}

			if typeFamily == u.TypeFamilyRELATIONSHIP {
				result += newline + "" + u.PascalToCamel(field.Name()) + " *g.Relationship"
			} else {
				result += newline + "" + u.PascalToCamel(field.Name()) + " *g." + getFieldForType(typeFamily)
			}
		}
	}

	return
}

func getFieldForType(typeFamily u.TypeFamily) string {
	switch typeFamily {
	case u.TypeFamilyBOOL:
		return "BoolField"
	case u.TypeFamilySTRING:
		return "StringField"
	case u.TypeFamilyINT:
		return "IntField"
	case u.TypeFamilyBIGINT:
		return "BigIntField"
	case u.TypeFamilyREAL:
		return "RealField"
	case u.TypeFamilyDOUBLE:
		return "DoubleField"
	case u.TypeFamilyDATE:
		return "DateField"
	case u.TypeFamilyENUM:
		return "EnumField"
	default:
		return typeFamily.String()
	}
}

func buildPropInit(classUtils IClassUtils, context *classGenContext) string {
	// the class as a variable
	className := u.PascalToCamel(string(classUtils.getClass()))

	// dealing with the class initialisation
	classInit := "newClass := &" + className + classNAMExSUFFIX + "{%s: %s}"
	superClassDecl := "IBusinessObjectClass"
	superClassValue := "g.NewClass()"
	if context.superType.Equals(typeURLxQUERYxOBJECT) {
		superClassDecl = "IURLQueryParamsClass"
		superClassValue = "g.NewURLQueryParamsClass()"
	} else if !context.superType.Equals(typeBUSINESSxOBJECT) {
		superClassDecl = u.PascalToCamel(context.superType.Name()) + classNAMExSUFFIX
		superClassValue = "*new" + context.superType.Name() + "Class()"
	}
	classInit = fmt.Sprintf(classInit, superClassDecl, superClassValue)

	// now adding the lines for the propertiess
	propLines := []string{classInit}

	// valueing each class property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		propLine := "newClass." + u.PascalToCamel(propName) + " = "

		multiple := "false"
		if propInfo.multiple {
			multiple = "true"
		}

		if propInfo.propType == u.TypeFamilyRELATIONSHIP {
			propLine += fmt.Sprintf("g.NewRelationship(%s, \"%s\", %s, %s)",
				"newClass", propName, multiple, u.PascalToCamel(propInfo.targetType))
		} else {
			if propInfo.propType == u.TypeFamilyENUM {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s, %s)",
					getFieldForType(propInfo.propType), "newClass", propName, multiple, "\""+propInfo.targetType+"\"")
			} else {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s)",
					getFieldForType(propInfo.propType), "newClass", propName, multiple)
			}
		}

		propLines = append(propLines, propLine+"")
	}

	// assembling the whole paragraph
	return strings.Join(propLines, newline)
}

func buildAccessors(classUtils IClassUtils, context *classGenContext) string {
	accessors := []string{}

	// generating 1 accessor per
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		owner := u.PascalToCamel(string(classUtils.getClass()))
		ownerShort := owner[:1]
		accType := getFieldForType(propInfo.propType)
		if propInfo.propType == u.TypeFamilyRELATIONSHIP {
			accType = "Relationship"
		}
		accessor := fmt.Sprintf("func (%s *%sClass) %s() *g.%s {"+
			newline+"return %s.%s"+
			newline+"}",
			ownerShort, owner, propName, accType,
			ownerShort, u.PascalToCamel(propName),
		)

		accessors = append(accessors, accessor)
	}

	return strings.Join(accessors, newline+newline)
}
