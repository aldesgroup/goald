// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the model files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
	"github.com/aldesgroup/goald/features/utils"
)

const modelTEMPLATE = `// Generated file, do not edit!
package model

import (
	"sync"

	g "github.com/aldesgroup/goald"$$imports$$
)

// static, reflect-free access to the definition of the $$Upper$$ model
type $$Upper$$Model struct {
$$propdecl$$
}

// this is the main way to refer to the $$Upper$$ model in the applicative code
func $$Upper$$() *$$Upper$$Model {
	return $$lower$$
}

// internal variables
var $$lower$$ *$$Upper$$Model
var $$lower$$Once sync.Once

// fully describing each of this model's properties & relationships
func New$$Upper$$Model() *$$Upper$$Model {
	$$propinit$$

	return thisModel
}

// making sure the $$Upper$$ model exists at app startup
func init() {
	$$lower$$Once.Do(func() {
		$$lower$$ = New$$Upper$$Model()
	})

	// this helps dynamically access to the $$Upper$$ model
	g.RegisterModel("$$Upper$$", $$lower$$)
}

// accessing all the $$Upper$$ model's properties and relationships

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
	core.EnsureDir(srcdir, includePATH)
	for _, includeDirEntry := range core.EnsureReadDir(srcdir, includePATH) {
		// if it's not a directory, we skip it
		if !includeDirEntry.IsDir() || includeDirEntry.Name() == dbFOLDERNAME {
			continue
		}

		// where the model files will be generated
		modelDir := path.Join(srcdir, includePATH, includeDirEntry.Name(), modelFOLDERxNAME)

		// we'll gather all the existing model files
		existingModelFiles := map[utils.ModelName]*modelFile{}

		// so, let's read the model folders
		for _, modelEntry := range core.EnsureReadDir(modelDir) {
			modelEntryInfo, errInfo := modelEntry.Info()
			core.PanicMsgIfErr(errInfo, "Could not read info for file '%s'", modelEntry.Name())
			if strings.HasSuffix(modelEntry.Name(), modelFILExSUFFIX) {
				modelName := utils.ModelName(core.KebabToPascal(modelEntry.Name()[:len(modelEntry.Name())-modelFILExSUFFIXxLEN]))
				existingModelFiles[modelName] = &modelFile{
					modTime:  modelEntryInfo.ModTime(),
					filename: modelEntry.Name(),
				}
			}
		}

		// let's see what we have in terms of business objects
		for name, source := range sourceRegistry.items {
			// considering only the business objects of THIS module
			// and no interface (at least for now)
			if source.isFromDir(includeDirEntry.Name()) && !source.isInterface() {
				// do we need to regen the model file?
				if existingModel := existingModelFiles[name]; regen ||
					existingModel == nil || existingModel.modTime.Before(source.getLastBOMod()) {
					// generating the missing or outdated model
					thisServer.generateOneModel(modelDir, source)

					// the code has changed
					codeChanged = true
				}

				// flagging this business object model as NOT unneeded (i.e. needed)
				delete(existingModelFiles, name)
			}
		}

		// removing the unneeded model files
		for _, unneededModel := range existingModelFiles {
			thisServer.Info(fmt.Sprintf("removing %s", unneededModel.filename))
			if errRem := os.Remove(path.Join(modelDir, unneededModel.filename)); errRem != nil {
				core.PanicMsgIfErr(errRem, "Could not delete model file '%s'", unneededModel.filename)
			}
		}
	}

	return
}

type modelGenerationContext struct {
	superType     reflection.GoaldType
	propertyNames []string
	propertiesMap map[string]modelGenPropertyInfo
}

type modelGenPropertyInfo struct {
	propType   propertyType
	multiple   bool
	targetType string
}

func (thisServer *server) generateOneModel(modelDir string, source IBusinessObjectModelSource) {
	// starting to build the file content, with the same context
	context := &modelGenerationContext{propertiesMap: map[string]modelGenPropertyInfo{}}

	// trivial filling of the template
	content := strings.ReplaceAll(modelTEMPLATE, "$$Upper$$", string(source.getName()))
	content = strings.ReplaceAll(content, "$$lower$$", core.PascalToCamel(string(source.getName())))

	// declaring the properties of the model, wether they are fields or relationships
	imports := map[string]string{}
	content = strings.Replace(content, "$$propdecl$$", buildPropDecl(source, context, imports), 1)

	// valueing the properties
	content = strings.Replace(content, "$$propinit$$", buildPropInit(source, context, imports), 1)

	// building the accessors to the properties
	content = strings.Replace(content, "$$accessors$$", buildAccessors(source, context), 1)

	// building the imports section	if len(imports) > 0 {
	importLines := core.GetSortedValues(imports)
	if len(importLines) > 0 {
		content = strings.Replace(content, "$$imports$$", newline+strings.Join(importLines, newline), 1)
	} else {
		content = strings.Replace(content, "$$imports$$", "", 1)
	}

	// writing to file
	core.WriteToFile(content, modelDir, core.PascalToKebab(string(source.getName()))+modelFILExSUFFIX)

	thisServer.Info(fmt.Sprintf("(Re-)generated model %s", source.getName()))
}

// this function helps declare 1 property (field or relationship) in the declaration of the model type
func buildPropDecl(source IBusinessObjectModelSource, context *modelGenerationContext, imports map[string]string) (result string) {
	// getting the object's type
	bObjType := reflection.TypeOf(source.NewObject(), true)

	// the very first property, field #0, MUST be the business object's super model
	superModelField := bObjType.Field(0)
	if !superModelField.IsAnonymous() || !reflection.PointerTo(superModelField.Type()).Implements(typeIxBUSINESSxOBJECT) {
		core.PanicMsg("%s: this object's first property should be the BO it inherits from, i.e."+
			"goald.BusinessObject, or one of its descendants", source.getName())
	}

	if context.superType = superModelField.Type(); context.superType.Equals(typeBUSINESSxOBJECT) {
		result += "g.IBusinessObjectModel"
	} else if context.superType.Equals(typeURLxQUERYxOBJECT) {
		result += "g.IURLQueryParamsModel"
	} else {
		result += "" + getImportPkg(imports, source, superModelField.Type().Name()) +
			superModelField.Type().Name() + modelNAMExSUFFIX
	}

	// browsing the entity's properties
	for fieldNum := 1; fieldNum < bObjType.NumField(); fieldNum++ {
		// getting the current field
		field := bObjType.Field(fieldNum)

		// detecting its type and multiplicity
		propertyType, multiple := detectPropertyType(field, typeIxBUSINESSxOBJECT, typeIxENUM)

		// adding to the context, and the model file content
		if propertyType != propertyTypeUNKNOWN {
			context.propertyNames = append(context.propertyNames, field.Name()) // we're keeping the original order

			targetType := ""                                    // makes no sense for basic BO fields...
			if propertyType == propertyTypeRELATIONSHIPxMONOM { // ... but it does for relationships
				entityType := core.IfThenElse(multiple, field.Type().Elem(), field.Type()) // e.g. *Object or []Object -> type Object
				targetType = entityType.Elem().Name()
			} else if propertyType == propertyTypeENUM { // or enums.
				targetType = field.Type().String()
			}

			// keeping track of the property's characteristics - this will be of use in the init function of the Model object
			context.propertiesMap[field.Name()] = modelGenPropertyInfo{propertyType, multiple, targetType}

			// writing out the property's declaration inside the Model object it belong to
			if propertyType.IsRelationship() {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g.Relationship"
			} else {
				result += newline + "" + core.PascalToCamel(field.Name()) + " *g." + getFieldForType(propertyType)
			}
		}
	}

	return
}

func getImportPackageLine(source IBusinessObjectModelSource) string {
	return reflection.TypeOf(source.NewObject(), true).PkgPath() // e.g. github.com/aldesgroup/project/group/packagename
}

func getImportModelLine(source IBusinessObjectModelSource) string {
	targetObjGoModule := core.Before(getImportPackageLine(source), source.getSrcPath())                     // e.g. github.com/aldesgroup/project/
	return fmt.Sprintf("%[1]s_model \"%[2]s_include/%[1]s/model\"", source.getPackage(), targetObjGoModule) // e.g. packagename_model "github.com/aldesgroup/project/_include/packagename"
}

func getFieldForType(propertyType propertyType) string {
	switch propertyType {
	case propertyTypeBOOL:
		return "BoolField"
	case propertyTypeSTRING:
		return "StringField"
	case propertyTypeINT:
		return "IntField"
	case propertyTypeBIGINT:
		return "BigIntField"
	case propertyTypeREAL:
		return "RealField"
	case propertyTypeDOUBLE:
		return "DoubleField"
	case propertyTypeDATE:
		return "DateField"
	case propertyTypeENUM:
		return "EnumField"
	default:
		return propertyType.String()
	}
}

// This function builds the line that helps initialise a model instance, for 1 property
func buildPropInit(source IBusinessObjectModelSource, context *modelGenerationContext, imports map[string]string) string {
	// dealing with the model initialisation
	modelInit := "thisModel := &" + string(source.getName()) + modelNAMExSUFFIX + "{%s: %s}"
	superModelDecl := "IBusinessObjectModel"
	superModelValue := "g.NewBusinessObjectModel()"
	if context.superType.Equals(typeURLxQUERYxOBJECT) {
		superModelDecl = "IURLQueryParamsModel"
		superModelValue = "g.NewURLQueryParamsModel()"
	} else if !context.superType.Equals(typeBUSINESSxOBJECT) {
		superModelDecl = context.superType.Name() + modelNAMExSUFFIX
		superModelValue = "*" + getImportPkg(imports, source, context.superType.Name()) + "New" + context.superType.Name() + "Model()"
	}
	modelInit = fmt.Sprintf(modelInit, superModelDecl, superModelValue)

	// now adding the lines for the propertiess
	propLines := []string{modelInit}

	// valueing each model property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		propLine := "thisModel." + core.PascalToCamel(propName) + " = "

		multiple := "false"
		if propInfo.multiple {
			multiple = "true"
		}

		if propInfo.propType == propertyTypeRELATIONSHIPxMONOM {
			propLine += fmt.Sprintf("g.AddRelationship(%s, \"%s\", \"%s\", %s, \"%s\")",
				"thisModel", source.getName(), propName, multiple, propInfo.targetType)
		} else if propInfo.propType == propertyTypeRELATIONSHIPxPOLYM {
			propLine += fmt.Sprintf("g.AddPolyRelationship(%s, \"%s\", \"%s\", %s)",
				"thisModel", source.getName(), propName, multiple)
		} else {
			if propInfo.propType == propertyTypeENUM {
				propLine += fmt.Sprintf("g.Add%s(%s, \"%s\", \"%s\", %s, %s)",
					getFieldForType(propInfo.propType), "thisModel", source.getName(), propName, multiple, "\""+propInfo.targetType+"\"")
			} else {
				propLine += fmt.Sprintf("g.Add%s(%s, \"%s\", \"%s\", %s)",
					getFieldForType(propInfo.propType), "thisModel", source.getName(), propName, multiple)
			}
		}

		propLines = append(propLines, propLine+"")
	}

	// assembling the whole paragraph
	return strings.Join(propLines, newline)
}

func getImportPkg(imports map[string]string, fromSource IBusinessObjectModelSource, forSourceName string) string {
	sourceName := utils.ModelName(forSourceName)
	sourceObj := sourceRegistry.items[sourceName]
	sourcePkg := sourceObj.getPackage()
	sourceMod := sourceObj.getModule()
	superImport := ""
	if sourceMod != getCurrentSourceModuleName() || sourcePkg != fromSource.getPackage() {
		if imports[sourceObj.getPackage()] == "" {
			imports[sourceObj.getPackage()] = getImportModelLine(sourceObj) // the needed import
		}
		superImport = sourcePkg + "_model."
	}

	return superImport
}

// This function builds an access for a property (field or relationship)
func buildAccessors(source IBusinessObjectModelSource, context *modelGenerationContext) string {
	accessors := []string{}

	// generating 1 accessor per property
	for _, propName := range context.propertyNames {
		propInfo := context.propertiesMap[propName]
		ownerName := source.getName()
		ownerShort := ownerName[:1]
		accType := getFieldForType(propInfo.propType)
		if propInfo.propType.IsRelationship() {
			accType = "Relationship"
		}
		accessor := fmt.Sprintf("func (%s *%sModel) %s() *g.%s {"+
			newline+"return %s.%s"+
			newline+"}",
			ownerShort, ownerName, propName, accType,
			ownerShort, core.PascalToCamel(propName),
		)

		accessors = append(accessors, accessor)
	}

	return strings.Join(accessors, newline+newline)
}
