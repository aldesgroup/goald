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
)

const sourceFILExSUFFIX = "__.go"
const sourceREGISTRYxNAME = "1_registry.go"

func (thisServer *server) generateObjectRegistries(srcdir, currentPath string, cache map[string]bool) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// reading the current directory
	dirEntries, errDir := os.ReadDir(readingPath)
	panicErrf(errDir, "could not read '%s'", readingPath)

	// we'll gather the business object infos here
	entriesInCode := map[string]*businessObjectEntry{}

	// going through the resources found withing the current directory
	for _, entry := range dirEntries {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" {
				// found another directory, let's dive deeper!
				thisServer.generateObjectRegistries(srcdir, path.Join(currentPath, entry.Name()), cache)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry for the egustry, from the current file
				bObjEntry := getEntryFromFile(path.Join(readingPath, entry.Name()))

				// checking the biz obj / file naming
				if PascalToSnake(bObjEntry.name)+sourceFILExSUFFIX != entry.Name() {
					panicf("The business object's name should be the file name Pascal-cased, i.e. we should have: "+
						"%s in file %s, "+
						"or %s in file %s",
						bObjEntry.name, PascalToSnake(bObjEntry.name)+sourceFILExSUFFIX,
						SnakeToPascal(strings.Replace(entry.Name(), sourceFILExSUFFIX, "", 1)), entry.Name(),
					)
				}

				// checking the unicity of each biz obj name
				if cache[bObjEntry.name] {
					panicf("We can't have 2 business objects with the same name '%s'."+
						" This would lead to the same REST path. You have to rename one.", bObjEntry.name)
				} else {
					entriesInCode[bObjEntry.name] = bObjEntry
					cache[bObjEntry.name] = true
				}
			}
		}
	}

	// do we need to regenerate the object registry at the current path?
	needRegen := false

	// let's check the current entries, the ones coded right now
	for entryName, entryIncode := range entriesInCode {
		entryInRegistry := boRegistry.content[entryName]
		if entryInRegistry == nil {
			log.Printf("Business object '%s' has appeared since the last generation!", entryName)
			needRegen = true

			break
		} else if entryInRegistry.lastMod.Before(entryIncode.lastMod) {
			log.Printf("Business object '%s' has changed since the last generation!", entryName)
			needRegen = true

			break
		}
	}

	// if we're not doing regen because of added or changed biz objs,
	// maybe we have to because of deleted ones!
	if !needRegen {
		for entryName := range boRegistry.folders[path.Join(currentMod(), currentPath)] {
			if entriesInCode[entryName] == nil {
				needRegen = true
				log.Printf("Business object '%s' has disappeard since the last generation!", entryName)

				break
			}
		}
	}

	// now let's write the registry file, if needed
	if needRegen && len(entriesInCode) > 0 {
		writeRegistryFile(srcdir, currentPath, entriesInCode)
	} else if filename := path.Join(srcdir, currentPath, sourceREGISTRYxNAME); len(entriesInCode) == 0 && FileExists(filename) {
		panicErrf(os.Remove(filename), "Could not delete '%s'", filename)
	}
}

const registryFileTemplate = `// Generated file, do not edit!
package %s

import (
	g "git-ext.aldes.com/j.wan/arch-poc/goald"
)

func init() {
%s
}
`

func writeRegistryFile(srcdir, currentPath string, bizObjEntries map[string]*businessObjectEntry) {
	// creating the file
	fileName := path.Join(srcdir, currentPath, sourceREGISTRYxNAME)

	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		panicf("Could not create file %s; cause: %s", fileName, errCreate)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Fatalf("Could not properly close file %s; cause: %s", fileName, errClose)
		}
	}()

	// gathering the biz objs in order
	// and the switch / case lines
	registrationLines := []string{}
	for _, bizObjEntry := range GetSortedValues[string, *businessObjectEntry](bizObjEntries) {
		registrationLines = append(registrationLines, fmt.Sprintf("%sg.Register(&%s{}, \"%s\")", "\t",
			bizObjEntry.name, bizObjEntry.lastMod.Add(oneSECOND).Format(dateFormatSECONDS)))
	}

	// writing to the file
	_, errWrite := file.WriteString(fmt.Sprintf(registryFileTemplate, path.Base(currentPath),
		strings.Join(registrationLines, newline)))
	panicErrf(errWrite, "Error while trying to (re)generate '%s'", fileName)
}

func getEntryFromFile(filename string) (entry *businessObjectEntry) {
	stat, errStat := os.Stat(filename)
	panicErrf(errStat, "Could not check the modification time for file '%s'", filename)

	// parsing the file to get the AST (Abstract Syntax Tree)
	file, errParse := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	panicErrf(errParse, "Error while parsing '%s'", filename)

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
								name:    typeSpec.Name.Name,
								lastMod: stat.ModTime(),
							}
						} else {
							panicf("More than one struct declared in the BusinessObject file '%s'!", filename)
						}
					}
				}
			}
		}
	}

	return
}

func currentMod() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panicf("Could not read the build infos!")
	}

	return bi.Main.Path
}
