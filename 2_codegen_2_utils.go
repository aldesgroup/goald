// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the class files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aldesgroup/goald/features/utils"
)

const utilsTEMPLATE = `// Generated file, do not edit!
package $$package$$

import (
	g "github.com/aldesgroup/goald"$$otherimports$$
)

// getting a property's value as a string, without using reflection
func (this$$Upper$$ *$$Upper$$) GetValueAsString(propertyName string) string {
	switch propertyName {
$$getcases$$
	default:
		return "unknown property: " + propertyName
	}
}

// setting a property's value with a given string value, without using reflection
func (this$$Upper$$ *$$Upper$$) SetValueAsString(propertyName string, valueAsString string) error {
	switch propertyName {
$$setcases$$
	}

	return g.Error("Unknown property: %T.%s", this$$Upper$$, propertyName)
}
`

const utilsFILExSUFFIX = "__utils.go"

func (thisServer *server) generateObjectUtils(srcdir, currentPath string, regen bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// reading the current directory
	dirEntries, errDir := os.ReadDir(readingPath)
	utils.PanicErrf(errDir, "could not read '%s'", readingPath)

	// going through the resources found withing the current directory
	for _, entry := range dirEntries {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				thisServer.generateObjectUtils(srcdir, path.Join(currentPath, entry.Name()), regen)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry within this file, then the registred entry in the code
				bObjEntryFile := getEntryFromFile(srcdir, currentPath, entry.Name())
				bObjEntry := boRegistry.content[bObjEntryFile.name]

				// the corresponding utils file, if it exist
				utilsFilename := strings.Replace(entry.Name(), "__.go", utilsFILExSUFFIX, 1)

				// generating the utils file, if not existing yet, or too old
				if regen || !utils.FileExists(utilsFilename) || utils.EnsureModTime(utilsFilename).Before(bObjEntry.lastMod) {
					generateObjectUtilsForEntry(srcdir, bObjEntry, utilsFilename)
				}
			}
		}
	}
}

func generateObjectUtilsForEntry(srcdir string, bObjEntry *businessObjectEntry, filename string) {
	content := strings.ReplaceAll(utilsTEMPLATE, "$$package$$", path.Base(bObjEntry.srcPath))
	content = strings.ReplaceAll(content, "$$Upper$$", bObjEntry.name)

	getCases := []string{}
	setCases := []string{}

	// need for some imports
	var importsMap = map[string]bool{}
	var importStrconv, importTime bool

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for fieldNum := 1; fieldNum < bObjEntry.bObjType.NumField(); fieldNum++ {
		// getting the current field
		field := bObjEntry.bObjType.Field(fieldNum)

		// detecting its type and multiplicity
		propType, multiple := getPropertyType(field)

		// adding to the context, and the class file content
		if propType != PropertyTypeUNKNOWN && propType != PropertyTypeRELATIONSHIP {
			if !multiple { // not handling this for now
				// case init
				getCase := fmt.Sprintf("\tcase \"%s\":", field.Name)
				setCase := getCase

				switch propType {
				case PropertyTypeBOOL:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.FormatBool(this%s.%s)", bObjEntry.name, field.Name)
					importStrconv = true
					setCase += newline + "\t\tvalue, errConv := strconv.ParseBool(valueAsString)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = value", bObjEntry.name, field.Name)

				case PropertyTypeSTRING:
					getCase += newline + fmt.Sprintf("\t\treturn this%s.%s", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = valueAsString", bObjEntry.name, field.Name)

				case PropertyTypeINT:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.Itoa(this%s.%s)", bObjEntry.name, field.Name)
					importStrconv = true
					setCase += newline + "\t\tvalue, errConv := strconv.Atoi(valueAsString)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = value", bObjEntry.name, field.Name)

				case PropertyTypeINT64:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.FormatInt(this%s.%s, 10)", bObjEntry.name, field.Name)
					importStrconv = true
					setCase += newline + "\t\tvalue, errConv := strconv.ParseInt(valueAsString, 10, 64)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = value", bObjEntry.name, field.Name)

				case PropertyTypeREAL32:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.FormatFloat(float64(this%s.%s), 'f', -1, 32)", bObjEntry.name, field.Name)
					importStrconv = true
					setCase += newline + "\t\tvalue, errConv := strconv.ParseFloat(valueAsString, 32)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = float32(value)", bObjEntry.name, field.Name)

				case PropertyTypeREAL64:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.FormatFloat(this%s.%s, 'f', -1, 64)", bObjEntry.name, field.Name)
					importStrconv = true
					setCase += newline + "\t\tvalue, errConv := strconv.ParseFloat(valueAsString, 64)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = value", bObjEntry.name, field.Name)

				case PropertyTypeDATE:
					getCase += newline + fmt.Sprintf("\t\treturn this%s.%s.Format(utils.RFC3339Milli)", bObjEntry.name, field.Name)
					importTime = true
					setCase += newline + "\t\tvalue, errConv := time.Parse(utils.RFC3339Milli, valueAsString)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = &value", bObjEntry.name, field.Name)

				case PropertyTypeENUM:
					getCase += newline + fmt.Sprintf("\t\treturn strconv.Itoa(this%s.%s.Val())", bObjEntry.name, field.Name)
					importStrconv = true
					enumTypeString := field.Type.Name() // e.g.: MyEnumType
					if enumPkg := field.Type.PkgPath(); enumPkg != bObjEntry.bObjType.PkgPath() {
						enumTypeString = field.Type.String() /// e.g.: thatpackage.MyEnumType
						importsMap[enumPkg] = true
					}
					setCase += newline + "\t\tintValue, errConv := strconv.Atoi(valueAsString)"
					setCase += newline + fmt.Sprintf("\t\tutils.PanicErrf(errConv, \"Could not set '%s.%s' to %%s\", valueAsString)", bObjEntry.name, field.Name)
					setCase += newline + fmt.Sprintf("\t\tthis%s.%s = (%s)(intValue)", bObjEntry.name, field.Name, enumTypeString)
					setCase += newline + fmt.Sprintf("\t\tutils.PanicIff(this%s.%s.String() == \"\", \"Could not set '%s.%s' to %%s since it's not a listed value\", valueAsString)",
						bObjEntry.name, field.Name, bObjEntry.name, field.Name)
				}

				// appending the case
				getCases = append(getCases, getCase)
				setCases = append(setCases, setCase)
			}
		}
	}

	// handling the imports
	content = strings.ReplaceAll(content, "$$getcases$$", strings.Join(getCases, newline))
	content = strings.ReplaceAll(content, "$$setcases$$", strings.Join(setCases, newline))

	if importStrconv {
		importsMap["strconv"] = true
	}
	if importTime {
		importsMap["time"] = true
	}
	if importStrconv || importTime {
		importsMap["github.com/aldesgroup/goald/features/utils"] = true
	}

	imports := ""
	if len(importsMap) > 0 {
		imports = newline + "\t\"" + strings.Join(utils.GetSortedKeys(importsMap), "\""+newline+"\t"+"\"") + "\""
	}
	content = strings.Replace(content, "$$otherimports$$", imports, 1)

	// write out the file
	utils.WriteToFile(content, srcdir, bObjEntry.srcPath, filename)
}
