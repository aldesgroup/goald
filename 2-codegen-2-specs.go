// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const specsTEMPLATE = `// Generated file, do not edit!
package specs

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

// static, reflect-free access to the definition of the $$Upper$$ specs
type $$lower$$Specs struct {
$$propdecl$$
}

// this is the main way to refer to the $$Upper$$ specs in the applicative code
func $$Upper$$() *$$lower$$Specs {
	return $$lower$$
}

// internal variables
var $$lower$$ *$$lower$$Specs
var $$lower$$Once sync.Once

// fully describing each of this class' properties & relationships
func new$$Upper$$Specs() *$$lower$$Specs {
	$$propinit$$

	return newSpecs
}

// making sure the $$Upper$$ specs exists at app startup
func init() {
	$$lower$$Once.Do(func() {
		$$lower$$ = new$$Upper$$Specs()
	})

	// this helps dynamically access to the $$Upper$$ specs
	g.RegisterSpecs("$$Upper$$", $$lower$$)
}

// accessing all the $$Upper$$ class' properties and relationships

$$accessors$$

`

const specsFOLDER = "_include/_specs"
const specsFILExSUFFIX = "--spc.go"
const specsFILExSUFFIXxLEN = len(specsFILExSUFFIX)
const specsNAMExSUFFIX = "Specs"
const newline = "\n"

func (thisServer *server) generateAllObjectSpecs(srcdir string, regen bool) (codeChanged bool) {
	// a type just used here
	type specsFile struct {
		modTime  time.Time
		filename string
	}

	// where the class files will be generated
	specsDir := core.EnsureDir(srcdir, specsFOLDER)

	// we'll gather all the existing class files
	existingSpecsFiles := map[className]*specsFile{}

	// so, let's read the class folder
	for _, specsEntry := range core.EnsureReadDir(specsDir) {
		specsEntryInfo, errInfo := specsEntry.Info()
		core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", specsEntry.Name())
		specsClassName := className(core.KebabToPascal(specsEntry.Name()[:len(specsEntry.Name())-specsFILExSUFFIXxLEN]))
		existingSpecsFiles[specsClassName] = &specsFile{
			modTime:  specsEntryInfo.ModTime(),
			filename: specsEntry.Name(),
		}
	}

	// let's see what we have in terms of business objects
	for name, class := range classRegistry.items {
		// considering only the business objects of THIS module
		// and no interface (at least for now)
		if class.getModule() == getCurrentModuleName() && !class.isInterface() {
			// do we need to regen the class file?
			if existingSpecs := existingSpecsFiles[name]; regen ||
				existingSpecs == nil || existingSpecs.modTime.Before(class.getLastBOMod()) {
				// generating the missing or outdated class
				generateOneSpecs(specsDir, class)

				// the code has changed
				codeChanged = true
			}

			// flagging this business object class as NOT unneeded (i.e. needed)
			delete(existingSpecsFiles, name)
		}
	}

	// removing the unneeded classes
	for _, unneededSpecs := range existingSpecsFiles {
		slog.Info(fmt.Sprintf("removing %s", unneededSpecs.filename))
		if errRem := os.Remove(path.Join(specsDir, unneededSpecs.filename)); errRem != nil {
			core.PanicMsgIfErr(errRem, "Could not delete class file '%s'", unneededSpecs.filename)
		}
	}

	return
}

type specsGenerationContext struct {
	superType     utils.GoaldType
	propertyNames []string
	propertiesMap map[string]classGenPropertyInfo
}

type classGenPropertyInfo struct {
	propType    utils.TypeFamily
	multiple    bool
	targetType  string
	targetTypes []string
}

func generateOneSpecs(specsDir string, class IClass) {
	// starting to build the file content, with the same context
	context := &specsGenerationContext{propertiesMap: map[string]classGenPropertyInfo{}}

	// trivial filling of the template
	clsName := string(class.getClassName())
	content := strings.ReplaceAll(specsTEMPLATE, "$$Upper$$", clsName)
	content = strings.ReplaceAll(content, "$$lower$$", core.PascalToCamel(clsName))

	// declaring the properties of the classe
	content = strings.Replace(content, "$$propdecl$$", buildPropDecl(class, context), 1)

	// valueing the properties
	content = strings.Replace(content, "$$propinit$$", buildPropInit(class, context), 1)

	// building the accessors to the properties
	content = strings.Replace(content, "$$accessors$$", buildAccessors(class, context), 1)

	// writing to file
	core.WriteToFile(content, specsDir, core.PascalToKebab(clsName)+specsFILExSUFFIX)

	slog.Info(fmt.Sprintf("(Re-)generated class %s", clsName))
}

// this function helps declare 1 property (field or relationship) in the declaration of the specs type
func buildPropDecl(class IClass, context *specsGenerationContext) (result string) {
	// getting the object's type
	bObjType := utils.TypeOf(class.NewObject(), true)

	// the very first property, field #0, MUST be the business object's super class
	superClassField := bObjType.Field(0)
	if !superClassField.IsAnonymous() || !utils.PointerTo(superClassField.Type()).Implements(typeIxBUSINESSxOBJECT) {
		core.PanicMsg("%s: this object's first property should be the BO it inherits from, i.e."+
			"goald.BusinessObject, or one of its descendants", class.getClassName())
	}

	if context.superType = superClassField.Type(); context.superType.Equals(typeBUSINESSxOBJECT) {
		result += "g.IBusinessObjectSpecs"
	} else if context.superType.Equals(typeURLxQUERYxOBJECT) {
		result += "g.IURLQueryParamsSpecs"
	} else {
		result += "" + core.PascalToCamel(superClassField.Type().Name()) + specsNAMExSUFFIX
	}

	// browsing the entity's properties
	for fieldNum := 1; fieldNum < bObjType.NumField(); fieldNum++ {
		// getting the current field
		field := bObjType.Field(fieldNum)

		// detecting its type and multiplicity
		typeFamily, multiple := utils.GetTypeFamily(field, typeIxBUSINESSxOBJECT, typeIxENUM)

		// adding to the context, and the class file content
		if typeFamily != utils.TypeFamilyUNKNOWN {
			context.propertyNames = append(context.propertyNames, field.Name()) // we're keeping the original order

			targetType := ""                                      // makes no sense for basic BO fields...
			var targetTypes []string                              // ... this even less...
			if typeFamily == utils.TypeFamilyRELATIONSHIPxMONOM { // ... but it does for relationships
				entityType := core.IfThenElse(multiple, field.Type().Elem(), field.Type())
				targetType = entityType.Elem().Name()
			} else if typeFamily == utils.TypeFamilyRELATIONSHIPxPOLYM { // ... but it does for relationships
				interfaceType := field.Type()
				if multiple {
					interfaceType = field.Type().Elem()
				}
				targetTypes = getImplementionsOfInterface(interfaceType)
			} else if typeFamily == utils.TypeFamilyENUM { // or enums.
				targetType = field.Type().String()
			}

			// keeping track of the property's characteristics - this will be of use in the init function of the Specs object
			context.propertiesMap[field.Name()] = classGenPropertyInfo{typeFamily, multiple, targetType, targetTypes}

			// writing out the property's declaration inside the Specs object it belong to
			if typeFamily.IsRelationship() {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g.Relationship"
			} else {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g." + getFieldForType(typeFamily)
			}
		}
	}

	return
}

func getFieldForType(typeFamily utils.TypeFamily) string {
	switch typeFamily {
	case utils.TypeFamilyBOOL:
		return "BoolField"
	case utils.TypeFamilySTRING:
		return "StringField"
	case utils.TypeFamilyINT:
		return "IntField"
	case utils.TypeFamilyBIGINT:
		return "BigIntField"
	case utils.TypeFamilyREAL:
		return "RealField"
	case utils.TypeFamilyDOUBLE:
		return "DoubleField"
	case utils.TypeFamilyDATE:
		return "DateField"
	case utils.TypeFamilyENUM:
		return "EnumField"
	default:
		return typeFamily.String()
	}
}

// This function builds the line that helps initialise a specs instance, for 1 property
func buildPropInit(class IClass, context *specsGenerationContext) string {
	// the class as a variable
	className := core.PascalToCamel(string(class.getClassName()))

	// dealing with the class initialisation
	specsInit := "newSpecs := &" + className + specsNAMExSUFFIX + "{%s: %s}"
	superSpecsDecl := "IBusinessObjectSpecs"
	superSpecsValue := "g.NewBusinessObjectSpecs()"
	if context.superType.Equals(typeURLxQUERYxOBJECT) {
		superSpecsDecl = "IURLQueryParamsSpecs"
		superSpecsValue = "g.NewURLQueryParamsSpecs()"
	} else if !context.superType.Equals(typeBUSINESSxOBJECT) {
		superSpecsDecl = core.PascalToCamel(context.superType.Name()) + specsNAMExSUFFIX
		superSpecsValue = "*new" + context.superType.Name() + "Specs()"
	}
	specsInit = fmt.Sprintf(specsInit, superSpecsDecl, superSpecsValue)

	// now adding the lines for the propertiess
	propLines := []string{specsInit}

	// valueing each class property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		propLine := "newSpecs." + core.PascalToCamel(propName) + " = "

		multiple := "false"
		if propInfo.multiple {
			multiple = "true"
		}

		if propInfo.propType == utils.TypeFamilyRELATIONSHIPxMONOM {
			propLine += fmt.Sprintf("g.NewRelationship(%s, \"%s\", %s, %s)",
				"newSpecs", propName, multiple, core.PascalToCamel(propInfo.targetType))
		} else if propInfo.propType == utils.TypeFamilyRELATIONSHIPxPOLYM {
			propLine += fmt.Sprintf("g.NewRelationship(%s, \"%s\", %s, %s)",
				"newSpecs", propName, multiple, strings.Join(core.MapFn(propInfo.targetTypes, core.PascalToCamel), ", "))
		} else {
			if propInfo.propType == utils.TypeFamilyENUM {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s, %s)",
					getFieldForType(propInfo.propType), "newSpecs", propName, multiple, "\""+propInfo.targetType+"\"")
			} else {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s)",
					getFieldForType(propInfo.propType), "newSpecs", propName, multiple)
			}
		}

		propLines = append(propLines, propLine+"")
	}

	// assembling the whole paragraph
	return strings.Join(propLines, newline)
}

// This function builds an access for a property (field or relationship)
func buildAccessors(class IClass, context *specsGenerationContext) string {
	accessors := []string{}

	// generating 1 accessor per
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		owner := core.PascalToCamel(string(class.getClassName()))
		ownerShort := owner[:1]
		accType := getFieldForType(propInfo.propType)
		if propInfo.propType.IsRelationship() {
			accType = "Relationship"
		}
		accessor := fmt.Sprintf("func (%s *%sSpecs) %s() *g.%s {"+
			newline+"return %s.%s"+
			newline+"}",
			ownerShort, owner, propName, accType,
			ownerShort, core.PascalToCamel(propName),
		)

		accessors = append(accessors, accessor)
	}

	return strings.Join(accessors, newline+newline)
}

var allInterfaceImplementations = map[string][]string{}

// this function finds all the implementations of a given interface
func getImplementionsOfInterface(interfaceType utils.GoaldType) []string {
	interfaceName := interfaceType.Name()

	// we may already have computed the answer...
	implementations := allInterfaceImplementations[interfaceName]

	// ...but maybe e haven't yet
	if implementations == nil {
		// browsing through all the non-interface classes to find the implementations
		for _, class := range classRegistry.items {
			if !class.isInterface() {
				if boType := utils.TypeOf(class.NewObject(), false); boType.Implements(interfaceType) {
					implementations = append(implementations, boType.Elem().Name())
				}
			}
		}

		// sorting, then caching for faster retrieval later
		slices.Sort(implementations)
		allInterfaceImplementations[interfaceName] = implementations
	}

	return implementations
}
