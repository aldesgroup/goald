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

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const modelTEMPLATE = `// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"$$imports$$
)

// static, reflect-free access to the definition of the $$Upper$$ model
type $$lower$$Model struct {
$$propdecl$$
}

// this is the main way to refer to the $$Upper$$ model in the applicative code
func $$Upper$$() *$$lower$$Model {
	return $$lower$$
}

// internal variables
var $$lower$$ *$$lower$$Model
var $$lower$$Once sync.Once

// fully describing each of this class' properties & relationships
func new$$Upper$$Model() *$$lower$$Model {
	$$propinit$$

	return newModel
}

// making sure the $$Upper$$ model exists at app startup
func init() {
	$$lower$$Once.Do(func() {
		$$lower$$ = new$$Upper$$Model()
	})

	// this helps dynamically access to the $$Upper$$ model
	g.RegisterModel("$$Upper$$", $$lower$$)
}

// accessing all the $$Upper$$ class' properties and relationships

$$accessors$$

`

const modelFOLDERxNAME = "model"
const modelFILExSUFFIX = "--mdl.go"
const modelFILExSUFFIXxLEN = len(modelFILExSUFFIX)
const modelNAMExSUFFIX = "Model"
const newline = "\n"

func (thisServer *server) generateAllObjectModels(srcdir string, regen bool) (codeChanged bool) {
	// a type just used here
	type modelFile struct {
		modTime  time.Time
		filename string
	}

	// iterating over each package for which we've already got a registry
	for _, modelDirEntry := range core.EnsureReadDir(srcdir, includePATH) {
		// if it's not a directory, we skip it
		if !modelDirEntry.IsDir() {
			continue
		}

		// where the model files will be generated
		modelDir := core.EnsureDir(srcdir, includePATH, modelDirEntry.Name(), modelFOLDERxNAME)

		// we'll gather all the existing class files
		existingModelFiles := map[className]*modelFile{}

		// so, let's read the model folders
		for _, modelEntry := range core.EnsureReadDir(modelDir) {
			modelEntryInfo, errInfo := modelEntry.Info()
			core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", modelEntry.Name())
			modelClassName := className(core.KebabToPascal(modelEntry.Name()[:len(modelEntry.Name())-modelFILExSUFFIXxLEN]))
			existingModelFiles[modelClassName] = &modelFile{
				modTime:  modelEntryInfo.ModTime(),
				filename: modelEntry.Name(),
			}
		}

		// let's see what we have in terms of business objects
		for name, class := range classRegistry.items {
			// considering only the business objects of THIS module
			// and no interface (at least for now)
			if class.getModule() == getCurrentModuleName() && class.getPackage() == modelDirEntry.Name() && !class.isInterface() {
				// do we need to regen the class file?
				if existingModel := existingModelFiles[name]; regen ||
					existingModel == nil || existingModel.modTime.Before(class.getLastBOMod()) {
					// generating the missing or outdated class
					generateOneModel(modelDir, class)

					// the code has changed
					codeChanged = true
				}

				// flagging this business object class as NOT unneeded (i.e. needed)
				delete(existingModelFiles, name)
			}
		}

		// removing the unneeded classes
		for _, unneededModel := range existingModelFiles {
			slog.Info(fmt.Sprintf("removing %s", unneededModel.filename))
			if errRem := os.Remove(path.Join(modelDir, unneededModel.filename)); errRem != nil {
				core.PanicMsgIfErr(errRem, "Could not delete class file '%s'", unneededModel.filename)
			}
		}

		// let's make the registry file import the model package, if it doesn't already
		if codeChanged {
			filename := path.Join(srcdir, includePATH, modelDirEntry.Name(), sourceREGISTRYxNAME)
			modelsImportPath := "_ \"" + path.Join(getCurrentModule(), includePATH, modelDirEntry.Name(), modelFOLDERxNAME) + "\""
			core.ReplaceInFile(filename, map[string]string{importPLACEHOLDER: modelsImportPath})
		}
	}

	return
}

type modelGenerationContext struct {
	superType     utils.GoaldType
	propertyNames []string
	propertiesMap map[string]classGenPropertyInfo
}

type classGenPropertyInfo struct {
	propType   utils.TypeFamily
	multiple   bool
	targetType string
	// targetTypes []string
}

func generateOneModel(modelDir string, class IClass) {
	// starting to build the file content, with the same context
	context := &modelGenerationContext{propertiesMap: map[string]classGenPropertyInfo{}}

	// trivial filling of the template
	clsName := string(class.getClassName())
	content := strings.ReplaceAll(modelTEMPLATE, "$$Upper$$", clsName)
	content = strings.ReplaceAll(content, "$$lower$$", core.PascalToCamel(clsName))

	// declaring the properties of the class, wether they are fields or relationships
	imports := map[string]string{}
	content = strings.Replace(content, "$$propdecl$$", buildPropDecl(class, context, imports), 1)

	// building the imports section	if len(imports) > 0 {
	importLines := core.GetSortedValues(imports)
	if len(importLines) > 0 {
		content = strings.Replace(content, "$$imports$$", newline+strings.Join(importLines, newline), 1)
	} else {
		content = strings.Replace(content, "$$imports$$", "", 1)
	}

	// valueing the properties
	content = strings.Replace(content, "$$propinit$$", buildPropInit(class, context), 1)

	// building the accessors to the properties
	content = strings.Replace(content, "$$accessors$$", buildAccessors(class, context), 1)

	// writing to file
	core.WriteToFile(content, modelDir, core.PascalToKebab(clsName)+modelFILExSUFFIX)

	slog.Info(fmt.Sprintf("(Re-)generated model %s", clsName))
}

// this function helps declare 1 property (field or relationship) in the declaration of the model type
func buildPropDecl(class IClass, context *modelGenerationContext, imports map[string]string) (result string) {
	// getting the object's type
	bObjType := utils.TypeOf(class.NewObject(), true)

	// the very first property, field #0, MUST be the business object's super class
	superClassField := bObjType.Field(0)
	if !superClassField.IsAnonymous() || !utils.PointerTo(superClassField.Type()).Implements(typeIxBUSINESSxOBJECT) {
		core.PanicMsg("%s: this object's first property should be the BO it inherits from, i.e."+
			"goald.BusinessObject, or one of its descendants", class.getClassName())
	}

	if context.superType = superClassField.Type(); context.superType.Equals(typeBUSINESSxOBJECT) {
		result += "g.IBusinessObjectModel"
	} else if context.superType.Equals(typeURLxQUERYxOBJECT) {
		result += "g.IURLQueryParamsModel"
	} else {
		result += "" + core.PascalToCamel(superClassField.Type().Name()) + modelNAMExSUFFIX
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
			if typeFamily == utils.TypeFamilyRELATIONSHIPxMONOM { // ... but it does for relationships
				entityType := core.IfThenElse(multiple, field.Type().Elem(), field.Type()) // e.g. *Object or []Object -> type Object
				targetClsName := className(entityType.Elem().Name())                       // e.g. Object
				targetClass := classForName((targetClsName))                               // e.g. ClassForObject
				targetObjPkg := targetClass.getPackage()                                   // e.g. packagename
				targetObjMod := targetClass.getModule()                                    // e.g. projectname
				if targetObjMod != getCurrentModuleName() || targetObjPkg != class.getPackage() {
					if imports[targetClass.getPackage()] == "" {
						imports[targetClass.getPackage()] = getImportLineForClass(targetClass) // the needed import
					}
					targetType = fmt.Sprintf("%s_model.%s()", targetObjPkg, targetClsName) // e.g. packagename_model.Object()
				} else {
					if class.getClassName() == targetClsName {
						// relationship to the same class => we need to use THIS model
						targetType = "newModel"
					} else {
						targetType = string(targetClsName)
					}
				}
			} else if typeFamily == utils.TypeFamilyENUM { // or enums.
				targetType = field.Type().String()
			}

			// keeping track of the property's characteristics - this will be of use in the init function of the Model object
			context.propertiesMap[field.Name()] = classGenPropertyInfo{typeFamily, multiple, targetType}

			// writing out the property's declaration inside the Model object it belong to
			if typeFamily.IsRelationship() {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g.Relationship"
			} else {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g." + getFieldForType(typeFamily)
			}
		}
	}

	return
}

func getImportLineForClass(targetClass IClass) string {
	targetObject := targetClass.NewObject()                                                                      // e.g. *Object (runtime instance)
	targetObjType := utils.TypeOf(targetObject, true)                                                            // e.g. Object (runtime type)
	targetObjFullPkg := targetObjType.PkgPath()                                                                  // e.g. github.com/aldesgroup/project/group/packagename
	targetObjGoModule := core.Before(targetObjFullPkg, targetClass.getSrcPath())                                 // e.g. github.com/aldesgroup/project/
	return fmt.Sprintf("%[1]s_model \"%[2]s_include/%[1]s/model\"", targetClass.getPackage(), targetObjGoModule) // e.g. packagename_model "github.com/aldesgroup/project/_include/packagename"
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

// This function builds the line that helps initialise a model instance, for 1 property
func buildPropInit(class IClass, context *modelGenerationContext) string {
	// the class as a variable
	className := core.PascalToCamel(string(class.getClassName()))

	// dealing with the class initialisation
	modelInit := "newModel := &" + className + modelNAMExSUFFIX + "{%s: %s}"
	superModelDecl := "IBusinessObjectModel"
	superModelValue := "g.NewBusinessObjectModel()"
	if context.superType.Equals(typeURLxQUERYxOBJECT) {
		superModelDecl = "IURLQueryParamsModel"
		superModelValue = "g.NewURLQueryParamsModel()"
	} else if !context.superType.Equals(typeBUSINESSxOBJECT) {
		superModelDecl = core.PascalToCamel(context.superType.Name()) + modelNAMExSUFFIX
		superModelValue = "*new" + context.superType.Name() + "Model()"
	}
	modelInit = fmt.Sprintf(modelInit, superModelDecl, superModelValue)

	// now adding the lines for the propertiess
	propLines := []string{modelInit}

	// valueing each class property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		propLine := "newModel." + core.PascalToCamel(propName) + " = "

		multiple := "false"
		if propInfo.multiple {
			multiple = "true"
		}

		if propInfo.propType == utils.TypeFamilyRELATIONSHIPxMONOM {
			propLine += fmt.Sprintf("g.NewRelationship(%s, \"%s\", %s, %s)",
				"newModel", propName, multiple, core.PascalToCamel(propInfo.targetType))
		} else if propInfo.propType == utils.TypeFamilyRELATIONSHIPxPOLYM {
			propLine += fmt.Sprintf("g.NewPolyRelationship(%s, \"%s\", %s)",
				"newModel", propName, multiple)
		} else {
			if propInfo.propType == utils.TypeFamilyENUM {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s, %s)",
					getFieldForType(propInfo.propType), "newModel", propName, multiple, "\""+propInfo.targetType+"\"")
			} else {
				propLine += fmt.Sprintf("g.New%s(%s, \"%s\", %s)",
					getFieldForType(propInfo.propType), "newModel", propName, multiple)
			}
		}

		propLines = append(propLines, propLine+"")
	}

	// assembling the whole paragraph
	return strings.Join(propLines, newline)
}

// This function builds an access for a property (field or relationship)
func buildAccessors(class IClass, context *modelGenerationContext) string {
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
		accessor := fmt.Sprintf("func (%s *%sModel) %s() *g.%s {"+
			newline+"return %s.%s"+
			newline+"}",
			ownerShort, owner, propName, accType,
			ownerShort, core.PascalToCamel(propName),
		)

		accessors = append(accessors, accessor)
	}

	return strings.Join(accessors, newline+newline)
}
