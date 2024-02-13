// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

const classTEMPLATE = `package class
 
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

const classFOLDER = "class"
const classFILExSUFFIX = "_cls.go"
const classFILExSUFFIXxLEN = len(classFILExSUFFIX)
const classNAMExSUFFIX = "Class"
const newline = "\n"

func (thisServer *server) generateObjectClasses(srcdir string) {
	// a type just used here
	type classFile struct {
		modTime  time.Time
		filename string
	}

	// checking the class folder exist, or creating it on the way
	classDir := path.Join(srcdir, classFOLDER)
	if !DirExists(classDir) {
		panicErrf(os.Mkdir(classDir, 0o777), "Could not create the class folder '%s'", classDir)
	}

	// we'll gather all the existing class files
	existingClassFiles := map[string]*classFile{}

	// so, let's read the class folder
	classEntries, errDir := os.ReadDir(classDir)
	panicErrf(errDir, "Could not read the class folder")
	for _, classEntry := range classEntries {
		classEntryInfo, errInfo := classEntry.Info()
		panicErrf(errInfo, "Could not read info for file '%s'", classEntry.Name())
		classEntryName := SnakeToPascal(classEntry.Name()[:len(classEntry.Name())-classFILExSUFFIXxLEN])
		existingClassFiles[classEntryName] = &classFile{
			modTime:  classEntryInfo.ModTime(),
			filename: classEntry.Name(),
		}
	}

	// let's see what we have in terms of business objects
	for name, bObjEntry := range boRegistry.content {
		// do we need to regen the class file?
		if existingClass := existingClassFiles[name]; existingClass == nil || existingClass.modTime.Before(bObjEntry.lastMod) ||
			true { // TODO remove
			// generating the missing or outdated class
			generateObjectClass(classDir, bObjEntry)
		}

		// flagging this business object class as NOT unneeded (i.e. needed)
		delete(existingClassFiles, name)
	}

	// removing the unneeded classes
	for _, unneededClass := range existingClassFiles {
		log.Println("removing " + unneededClass.filename)
		if errRem := os.Remove(path.Join(classDir, unneededClass.filename)); errRem != nil {
			panicErrf(errRem, "Could not delete class file '%s'", unneededClass.filename)
		}
	}
}

type classGenContext struct {
	superType     reflect.Type
	propertyNames []string
	propertiesMap map[string]*classGenPropertyInfo
}

type classGenPropertyInfo struct {
	propType   PropertyType
	multiple   bool
	targetType string
}

func generateObjectClass(classDir string, bObjEntry *businessObjectEntry) {
	// creating the file
	fileName := path.Join(classDir, PascalToSnake(bObjEntry.name)+classFILExSUFFIX)

	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		panicf("Could not create file %s; cause: %s", fileName, errCreate)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Fatalf("Could not properly close file %s; cause: %s", fileName, errClose)
		}
	}()

	// starting to build the file content, with the same context
	context := &classGenContext{propertiesMap: map[string]*classGenPropertyInfo{}}

	// trivial filling of the template
	content := strings.ReplaceAll(classTEMPLATE, "$$Upper$$", bObjEntry.name)
	content = strings.ReplaceAll(content, "$$lower$$", PascalToCamel(bObjEntry.name))

	// declaring the properties of the classe
	content = strings.Replace(content, "$$propdecl$$", buildPropDecl(bObjEntry, context), 1)

	// valueing the properties
	content = strings.Replace(content, "$$propinit$$", buildPropInit(bObjEntry, context), 1)

	// building the accessors to the properties
	content = strings.Replace(content, "$$accessors$$", buildAccessors(bObjEntry, context), 1)

	// writing to file
	if _, errWrite := file.WriteString(content); errWrite != nil {
		panicErrf(errWrite, "Could not write file '%s'", fileName)
	}
}

func buildPropDecl(bObjEntry *businessObjectEntry, context *classGenContext) (result string) {
	// the very first property, field #0, MUST be the business object's super class
	superClassField := bObjEntry.bObjType.Field(0)
	if !superClassField.Anonymous || !reflect.PointerTo(superClassField.Type).Implements(typeIxBUSINESSxOBJECT) {
		panicf("%s: this object's first property should be the BO it inherits from, i.e."+
			"goald.BusinessObject, or one of its descendants", bObjEntry.name)
	}

	if context.superType = superClassField.Type; context.superType == typeBUSINESSxOBJECT {
		result += "\tg.IBusinessObjectClass"
	} else {
		result += PascalToCamel(superClassField.Type.Name()) + classNAMExSUFFIX
	}

	// browsing the entity's properties
	for fieldNum := 1; fieldNum < bObjEntry.bObjType.NumField(); fieldNum++ {
		// getting the current field
		field := bObjEntry.bObjType.Field(fieldNum)

		// detecting its type and multiplicity
		propType, multiple := getPropertyType(field)

		// adding to the context, and the class file content
		if propType != PropertyTypeUNKNOWN {
			context.propertyNames = append(context.propertyNames, field.Name) // we're keeping the original order

			targetType := ""                          // makes no sense for BO fields...
			if propType == PropertyTypeRELATIONSHIP { // ... but it does for relationships
				entityType := field.Type
				if multiple {
					entityType = field.Type.Elem()
				}
				targetType = entityType.Elem().Name()
			}

			context.propertiesMap[field.Name] = &classGenPropertyInfo{propType, multiple, targetType}

			if propType == PropertyTypeRELATIONSHIP {
				result += newline + PascalToCamel(field.Name) + " *g.Relationship"
			} else {
				result += newline + PascalToCamel(field.Name) + " *g." + getFieldForType(propType)
			}
		}
	}

	return
}

func getFieldForType(propType PropertyType) string {
	switch propType {
	case PropertyTypeBOOL:
		return "BoolField"
	case PropertyTypeSTRING:
		return "StringField"
	case PropertyTypeINT:
		return "IntField"
	case PropertyTypeINT64:
		return "Int64Field"
	case PropertyTypeREAL32:
		return "Real32Field"
	case PropertyTypeREAL64:
		return "Real64Field"
	case PropertyTypeDATE:
		return "DateField"
	case PropertyTypeENUM:
		return "EnumField"
	default:
		return propType.String()
	}
}

func buildPropInit(bObjEntry *businessObjectEntry, context *classGenContext) string {
	// the class as a variable
	className_ := PascalToCamel(bObjEntry.name)

	// dealing with the class initialisation
	classInit := "newClass := &" + className_ + classNAMExSUFFIX + "{%s: %s}"
	superClassDecl := "IBusinessObjectClass"
	superClassValue := "g.NewClass()"
	if context.superType != typeBUSINESSxOBJECT {
		superClassDecl = PascalToCamel(context.superType.Name()) + classNAMExSUFFIX
		superClassValue = "*new" + context.superType.Name() + "Class()"
	}
	classInit = fmt.Sprintf(classInit, superClassDecl, superClassValue)

	// now adding the lines for the propertiess
	propLines := []string{classInit}

	// valueing each class property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		propLine := "newClass." + PascalToCamel(propName) + " = "

		multiple := "false"
		if propInfo.multiple {
			multiple = "true"
		}

		if propInfo.propType == PropertyTypeRELATIONSHIP {
			propLine += fmt.Sprintf("g.NewRelationship(%s, \"%s\", %s, %s)",
				"newClass", propName, multiple, PascalToCamel(propInfo.targetType))
		} else {
			propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s)",
				getFieldForType(propInfo.propType), "newClass", propName, multiple)
		}

		propLines = append(propLines, propLine)
	}

	// assembling the whole paragraph
	return strings.Join(propLines, newline)
}

func buildAccessors(bObjEntry *businessObjectEntry, context *classGenContext) string {
	accessors := []string{}

	// generating 1 accessor per
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		owner := PascalToCamel(bObjEntry.name)
		ownerShort := owner[:1]
		accType := getFieldForType(propInfo.propType)
		if propInfo.propType == PropertyTypeRELATIONSHIP {
			accType = "Relationship"
		}
		accessor := fmt.Sprintf("func (%s *%sClass) %s() *g.%s {"+
			newline+"return %s.%s"+
			newline+"}",
			ownerShort, owner, propName, accType,
			ownerShort, PascalToCamel(propName),
		)

		// func (p *projectClass) Name() *g.Field {
		// 	return p.name
		// }

		accessors = append(accessors, accessor)
	}

	return strings.Join(accessors, newline+newline)
}
