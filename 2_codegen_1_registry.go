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
)

const sourceFILExSUFFIX = "__.go"
const sourceREGISTRYxNAME = "main/registry.go"

// TODO to have only 1 file, we need to handle the imports, which can be tedious, in particular with nested packages

func (thisServer *server) generateObjectRegistry(srcdir, currentPath string, entriesInCode map[string]*businessObjectEntry, start time.Time) {
	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// reading the current directory
	dirEntries, errDir := os.ReadDir(readingPath)
	panicErrf(errDir, "could not read '%s'", readingPath)

	// going through the resources found withing the current directory
	for _, entry := range dirEntries {
		if entry.IsDir() {
			// not going into the vendor
			if entry.Name() != "vendor" {
				// found another directory, let's dive deeper!
				thisServer.generateObjectRegistry(srcdir, path.Join(currentPath, entry.Name()), entriesInCode, start)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry for the egustry, from the current file
				if bObjEntry := getEntryFromFile(srcdir, currentPath, entry.Name()); bObjEntry != nil {
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
					if entriesInCode[bObjEntry.name] != nil {
						panicf("We can't have 2 business objects with the same name '%s'."+
							" This would lead to the same REST path. You have to rename one.", bObjEntry.name)
					} else {
						entriesInCode[bObjEntry.name] = bObjEntry
					}
				}
			}
		}
	}

	// if we're at root here, this means we've browsed through all the code already,
	// and can now decide to (re-)generate the object registry - or not
	if currentPath == "." {
		writeRegistryFileIfNeeded(srcdir, entriesInCode, start)
	}
}

const registryFileTemplate = `// Generated file, do not edit!
package main

import (
	g "github.com/aldesgroup/goald"
%s
)

func init() {
%s
}
`

func writeRegistryFileIfNeeded(srcdir string, entriesInCode map[string]*businessObjectEntry, start time.Time) {
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
		for entryName := range boRegistry.content {
			if entriesInCode[entryName] == nil {
				needRegen = true
				log.Printf("Business object '%s' has disappeard since the last generation!", entryName)

				break
			}
		}
	}

	// now let's write the registry file, if needed, and if we're at root
	if needRegen {
		// gathering the biz objs in order
		registrationLines := []string{}

		// and the imports, but only once per import, hence the map
		imports := []string{}
		imported := map[string]bool{}

		// going over all the business object entries
		for _, bObjEntry := range GetSortedValues[string, *businessObjectEntry](entriesInCode) {
			// adding 1 registration line per business object
			registrationLines = append(registrationLines, fmt.Sprintf("%sg.Register(&%s.%s{}, \"%s\")", "\t",
				path.Base(bObjEntry.pkgPath), bObjEntry.name, bObjEntry.lastMod.Add(oneSECOND).Format(dateFormatSECONDS)))

			// adding the corresponding import
			if !imported[bObjEntry.pkgPath] {
				imports = append(imports, "\""+bObjEntry.pkgPath+"\"")
				imported[bObjEntry.pkgPath] = true
			}
		}
		// writing to the file

		WriteToFile(
			fmt.Sprintf(registryFileTemplate, strings.Join(imports, newline), strings.Join(registrationLines, newline)),
			srcdir, sourceREGISTRYxNAME,
		)

		println(fmt.Sprintf("BO registry generated in %s", time.Since(start)))
	}
}

func getEntryFromFile(srcdir, entryPath, entryName string) (entry *businessObjectEntry) {
	filename := path.Join(srcdir, entryPath, entryName)

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
								pkgPath: path.Join(getCurrentModule(), entryPath),
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

var currentModule string

func getCurrentModule() string {
	if currentModule == "" {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			panicf("Could not read the build infos!")
		}

		currentModule = bi.Main.Path
	}

	return currentModule
}
