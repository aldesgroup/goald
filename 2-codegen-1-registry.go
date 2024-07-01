// ------------------------------------------------------------------------------------------------
// Here is the code used to generate the BO registry files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

const sourceFILExSUFFIX = "--.go"
const sourceREGISTRYxNAME = "0-registry.go"

// TODO to have only 1 file, we need to handle the imports, which can be tedious, in particular with nested packages

func (thisServer *server) generateObjectRegistry(srcdir, currentPath string, _ bool,
	allEntriesInCodeSoFar map[className]*businessObjectEntry, regen bool) {
	// we want all the entities we find in the code to build 1 global registry
	entriesInCode := allEntriesInCodeSoFar

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
				thisServer.generateObjectRegistry(srcdir, path.Join(currentPath, entry.Name()), false, entriesInCode, regen)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry for the egustry, from the current file
				if bObjEntry := getEntryFromFile(srcdir, currentPath, entry.Name()); bObjEntry != nil {
					// checking the biz obj / file naming
					if expected := utils.PascalToKebab(string(bObjEntry.class)) + sourceFILExSUFFIX; expected != entry.Name() {
						utils.Panicf("The business object's name should be the file name Pascal-cased, i.e. we should have: "+
							"%s in file %s, "+
							"or %s in file %s",
							bObjEntry.class, expected,
							utils.KebabToPascal(strings.Replace(entry.Name(), sourceFILExSUFFIX, "", 1)), entry.Name(),
						)
					}

					// checking the unicity of each biz obj name
					if entriesInCode[bObjEntry.class] != nil {
						utils.Panicf("We can't have 2 business objects with the same name '%s'."+
							" This would lead to the same REST path. You have to rename one.", bObjEntry.class)
					} else {
						entriesInCode[bObjEntry.class] = bObjEntry
					}
				}
			}
		}
	}

	// if we're at root here, this means we've browsed through all the code already,
	// and can now decide to (re-)generate the object registry - or not
	if currentPath == "." {
		writeRegistryFileIfNeeded(srcdir, currentPath, false, entriesInCode, regen)
	}
}

const registryFileTemplate = `// Generated file, do not edit!
package %s

import (
	g "github.com/aldesgroup/goald"
%s
)

func init() {
%s
}
`

func writeRegistryFileIfNeeded(srcdir, currentPath string, isLibrary bool,
	entriesInCode map[className]*businessObjectEntry, regen bool) {
	// do we need to regenerate the object registry at the current path?
	needRegen := regen

	// let's check the current entries, the ones coded right now
	for entryName, entryInCode := range entriesInCode {
		entryInRegistry := boRegistry.content[entryName]
		if entryInRegistry == nil {
			log.Printf("Business object '%s' has appeared since the last generation!", entryName)
			needRegen = true

			break
		} else if entryInRegistry.lastMod.Before(entryInCode.lastMod) {
			log.Printf("Business object '%s' has changed since the last generation!", entryName)
			needRegen = true

			break
		}
	}

	// if we're not doing regen because of added or changed biz objs,
	// maybe we have to because of deleted ones!
	if !needRegen {
		for entryName, entryInRegistry := range boRegistry.content {
			// in library mode, we're only considering the entries of the current source path
			if isLibrary {
				if entryInRegistry.srcPath != currentPath {
					continue
				}
			}

			if entriesInCode[entryName] == nil {
				needRegen = true
				log.Printf("Business object '%s' has disappeard since the last generation!", entryName)

				break
			}
		}
	}

	// now let's write the registry file, if needed, and if we're at root
	if nbEntries := len(entriesInCode); nbEntries > 0 && needRegen {
		// gathering the biz objs in order
		registrationLines := []string{fmt.Sprintf("\tg.In(\"%s\")", getCurrentModuleName())}

		// and the imports, but only once per import, hence the map
		imports := []string{}
		imported := map[string]bool{}

		// going over all the business object entries
		for _, bObjEntry := range utils.GetSortedValues[className, *businessObjectEntry](entriesInCode) {
			if isLibrary {
				// adding 1 registration line per business object
				boPath := path.Base(bObjEntry.srcPath)
				registrationLines = append(registrationLines,
					fmt.Sprintf(
						"%sRegister(func() any { return &%s.%s{} }, \"%s\", \"%s\", func() any { return []*%s.%s{} })", "\t\t",
						boPath, bObjEntry.class, bObjEntry.srcPath, bObjEntry.lastMod.Add(time.Second).Format(time.RFC3339), boPath, bObjEntry.class),
				)
			} else {
				// adding 1 registration line per business object
				boPath := path.Base(bObjEntry.srcPath)
				registrationLines = append(registrationLines,
					fmt.Sprintf(
						"%sRegister(func() any { return &%s.%s{} }, \"%s\", \"%s\", func() any { return []*%s.%s{} })", "\t\t",
						boPath, bObjEntry.class, bObjEntry.srcPath, bObjEntry.lastMod.Add(time.Second).Format(time.RFC3339), boPath, bObjEntry.class),
				)

				// adding the corresponding import
				if !imported[bObjEntry.srcPath] {
					imports = append(imports, "\""+path.Join(getCurrentModule(), bObjEntry.srcPath)+"\"")
					imported[bObjEntry.srcPath] = true
				}
			}
		}

		// where to write the file?
		genPath := "main"
		if isLibrary {
			genPath = currentPath
		}

		// which package?
		pkgName := path.Base(genPath)

		// which file?
		filename := path.Join(srcdir, genPath, sourceREGISTRYxNAME)

		// which content?
		dot := "." + newline
		content := fmt.Sprintf(registryFileTemplate, pkgName, strings.Join(imports, dot), strings.Join(registrationLines, dot))

		// writing to the file
		utils.WriteToFile(content, filename)
	}
}

func getEntryFromFile(srcdir, currentPath, entryName string) (entry *businessObjectEntry) {
	filename := path.Join(srcdir, currentPath, entryName)

	stat, errStat := os.Stat(filename)
	utils.PanicErrf(errStat, "Could not check the modification time for file '%s'", filename)

	// parsing the file to get the AST (Abstract Syntax Tree)
	file, errParse := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	utils.PanicErrf(errParse, "Error while parsing '%s'", filename)

	// going through the declarations in the file
	for _, node := range file.Decls {
		// in particular the generic declarations - as opposed to functions or bad declarations
		switch genDecl := node.(type) {
		case *ast.GenDecl:
			// going through the "specs" in the current declaration
			for _, spec := range genDecl.Specs {
				// stopping for declarations of type "type"
				switch typeSpec := spec.(type) {
				case *ast.TypeSpec:
					// more precisely, stopping for "struct" declarations
					switch typeSpec.Type.(type) {
					case *ast.StructType:
						if entry == nil {
							entry = &businessObjectEntry{
								class:   className(typeSpec.Name.Name),
								lastMod: stat.ModTime(),
								srcPath: currentPath,
							}
						} else {
							utils.Panicf("More than one struct declared in the BusinessObject file '%s'!", filename)
						}
					}
				}
			}
		}
	}

	return
}

var (
	currentModule     string
	currentModuleName string
)

// returns this module's path, e.g. "github.com/aldesgroup/goald"
func getCurrentModule() string {
	if currentModule == "" {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			utils.Panicf("Could not read the build infos!")
		}

		currentModule = bi.Main.Path
	}

	return currentModule
}

// returns this module's name, e.g. "goald"
func getCurrentModuleName() string {
	if currentModuleName == "" {
		currentModuleName = path.Base(getCurrentModule())
	}

	return currentModuleName
}
