// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the BO models in the web app
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

var modelsDIRPATH = path.Join("src", "components", "models")

func (thisServer *server) generateWebAppModels(webdir string, regen bool) {
	// the enum files to generate
	enums := map[string]IEnum{}

	// scanning for BOs involved in endpoints used from the web app
	for _, ep := range restRegistry.endpoints {
		if ep.isFromWebApp() {
			// generating the model for the endpoint resource
			thisServer.generateWebAppModel(webdir, ep.getResourceClass(), enums, regen)

			// if the endpoint admits a BO as an input (body or URL params), then we also need the model in the webapp
			if inputOrParamsClass := ep.getInputOrParamsClass(); inputOrParamsClass != "" {
				thisServer.generateWebAppModel(webdir, inputOrParamsClass, enums, regen)
			}
		}
	}

	// enum files generation
	for enumType, enum := range enums {
		// TODO run in go routines
		thisServer.generateWebAppEnum(webdir, enumType, enum, regen)
	}
}

func (thisServer *server) generateWebAppModel(webdir string, clsName className, enums map[string]IEnum, regen bool) {
	boClass := classForName(clsName)
	clUtils := getClassUtils(boClass)

	filename := utils.PascalToCamel(string(clsName)) + ".ts"
	filepath := path.Join(webdir, modelsDIRPATH, filename)

	// TODO remove
	// if utils.FileExists(filepath) && utils.EnsureModTime(filepath).After(clUtils.getLastBOMod()) {
	// 	return // the file already exists and is older than our changes in the BO class file
	// }

	// getting the type of business object
	bObjType := utils.TypeOf(clUtils.NewObject(), true)
	boInstance := utils.ValueOf(clUtils.NewObject())

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range utils.GetSortedValues(boClass.base().fields) {
		// adding to the context, and the class file content
		if typeFamily := field.getTypeFamily(); typeFamily != utils.TypeFamilyUNKNOWN && typeFamily != utils.TypeFamilyRELATIONSHIP {
			// not handling multiple properties for now
			if fieldName := field.getName(); !field.isMultiple() {
				// is the field type a type alias, or a built-in type?
				// fieldTypeAlias := getNonBuiltInFieldType(bObjectType, fieldName, importsMap)

				switch typeFamily {
				case utils.TypeFamilyENUM:
					// flagging this enum types for code generation
					enumType := bObjType.FieldByName(fieldName).Type().Name()
					enum := boInstance.GetFieldValue(fieldName).(IEnum)
					enums[enumType] = enum
				}

				// TODO handle the rest
			}
		}
	}

	utils.WriteToFile(fmt.Sprintf("// Coucou @ %s", time.Now()), filepath)
}

func (thisServer *server) generateWebAppEnum(webdir string, enumType string, enum IEnum, regen bool) {
	filename := "_" + utils.PascalToCamel(enumType) + ".ts"
	filepath := path.Join(webdir, modelsDIRPATH, filename)

	content := ""
	allTypes := []string{}
	for _, enumVal := range utils.GetSortedKeys(enum.Values()) {
		enumLabel := enum.Values()[enumVal]
		enumName := makeEnumName(enumLabel) // TODO : take existing enum name if any; or maybe NOT if the labels in the code are only keys to be translated
		content += fmt.Sprintf("export const %s = %d;", enumName, enumVal) + newline
		allTypes = append(allTypes, fmt.Sprintf("    { value: %s, label: \"%s\" },", enumName, enumLabel))
	}

	content += "export const Options = [" + newline
	content += strings.Join(allTypes, newline) + newline
	content += "];"

	utils.WriteToFile(content, filepath)
}

func makeEnumName(label string) string {
	sanitized := ""
	for _, rune := range label {
		switch rune {
		case '-', '(', ')', '&', '.':
			sanitized += " "
		case 'ç':
			sanitized += "c"
		default:
			sanitized += string(rune)
		}
	}

	bits := []string{}
	for _, bit := range strings.Split(sanitized, " ") {
		if bit != "" {
			bits = append(bits, bit)
		}
	}

	return strings.ToUpper(strings.Join(bits, "_"))
}
