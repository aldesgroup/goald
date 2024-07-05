// ------------------------------------------------------------------------------------------------
// Here is the code used to generate:
// - the BO registry files, which allows to register the corresponding ClassUtils objects
// - all the ClassUtils files
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"log/slog"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

const sourceFILExSUFFIX = "--.go"
const sourceCLSUxSUFFIX = "--clsu.go"
const sourceREGISTRYxNAME = "0-registry.go"
const sourceCLASSxUTILSxDIR = "classutils"

// ------------------------------------------------------------------------------------------------
// Going over all the physical source code files and generating stuff along the way
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateObjectRegistry(srcdir, currentPath string, _ bool,
	allEntriesInCodeSoFar map[className]*classUtilsCore, regen bool) {
	// we want all the entities we find in the code to build 1 global registry
	allClsuCoresInCode := allEntriesInCodeSoFar

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
				thisServer.generateObjectRegistry(srcdir, path.Join(currentPath, entry.Name()), false, allClsuCoresInCode, regen)
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// getting the business object entry for the egustry, from the current file
				if clsuCore := getClassUtilsFromFile(srcdir, currentPath, entry.Name()); clsuCore != nil {
					// checking the biz obj / file naming
					if expected := utils.PascalToKebab(string(clsuCore.class)) + sourceFILExSUFFIX; expected != entry.Name() {
						utils.Panicf("The business object's name should be the file name Pascal-cased, i.e. we should have: "+
							"%s in file %s, "+
							"or %s in file %s",
							clsuCore.class, expected,
							utils.KebabToPascal(strings.Replace(entry.Name(), sourceFILExSUFFIX, "", 1)), entry.Name(),
						)
					}

					// checking the unicity of each biz obj name
					if allClsuCoresInCode[clsuCore.class] != nil {
						utils.Panicf("We can't have 2 business objects with the same name '%s'."+
							" This would lead to the same REST path. You have to rename one.", clsuCore.class)
					} else {
						// adding one more BO to our list
						allClsuCoresInCode[clsuCore.class] = clsuCore

						// generating the corresponding ClassUtils file, if it doesn't exist yet
						genClassUtilsFile(srcdir, clsuCore)
					}
				}
			}
		}
	}

	// if we're at root here, this means we've browsed through all the code already,
	// and can now decide to (re-)generate the object registry - or not
	if currentPath == "." {
		writeRegistryFileIfNeeded(srcdir, allClsuCoresInCode, regen)
	}
}

// ------------------------------------------------------------------------------------------------
// Writing the registry file
// ------------------------------------------------------------------------------------------------

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

func writeRegistryFileIfNeeded(srcdir string, allClsuCoresInCode map[className]*classUtilsCore, regen bool) {
	// do we need to regenerate the object registry at the current path?
	needRegen := regen

	// let's check the current class utilS, the ones coded right now
	for clsName, clsuCoreInCode := range allClsuCoresInCode {
		classUtilsInRegistry := classUtilsRegistry.content[clsName]
		if classUtilsInRegistry == nil {
			log.Printf("Business object '%s' has appeared since the last generation!", clsName)
			needRegen = true

			break
		} else if classUtilsInRegistry.getLastBOMod().Before(clsuCoreInCode.getLastBOMod()) {
			log.Printf("Business object '%s' has changed since the last generation!", clsName)
			needRegen = true

			break
		}
	}

	// if we're not doing regen because of added or changed biz objs,
	// maybe we have to because of deleted ones!
	if !needRegen {
		for clsName := range classUtilsRegistry.content {
			if allClsuCoresInCode[clsName] == nil {
				needRegen = true
				log.Printf("Business object '%s' has disappeard since the last generation!", clsName)

				break
			}
		}
	}

	// now let's write the registry file, if needed, and if we're at root
	if nbEntries := len(allClsuCoresInCode); nbEntries > 0 && needRegen {
		// gathering the biz objs in order
		registrationLines := []string{fmt.Sprintf("\tg.In(\"%s\")", getCurrentModuleName())}

		// and the imports, but only once per import, hence the map
		imports := []string{}
		imported := map[string]bool{}

		// going over all the class utils cores
		for _, clsuCore := range utils.GetSortedValues[className, *classUtilsCore](allClsuCoresInCode) {
			// adding 1 registration line per business object
			boPath := path.Base(clsuCore.srcPath)
			registrationLines = append(registrationLines,
				fmt.Sprintf(
					"%sRegister(%s.ClassUtilsFor%s(\"%s\", \"%s\"))", "\t\t",
					boPath, clsuCore.class, clsuCore.srcPath, clsuCore.getLastBOMod().Add(time.Second).Format(time.RFC3339)),
			)

			// adding the corresponding import
			if !imported[clsuCore.srcPath] {
				imports = append(imports, boPath+" \""+path.Join(getCurrentModule(), clsuCore.srcPath, sourceCLASSxUTILSxDIR)+"\"")
				imported[clsuCore.srcPath] = true
			}
		}

		// where to write the file?
		genPath := "main"

		// which package?
		pkgName := path.Base(genPath)

		// which file?
		filename := path.Join(srcdir, genPath, sourceREGISTRYxNAME)

		// which content?
		dot := "." + newline
		content := fmt.Sprintf(registryFileTemplate, pkgName, strings.Join(imports, newline), strings.Join(registrationLines, dot))

		// writing to the file
		utils.WriteToFile(content, filename)
	}
}

// ------------------------------------------------------------------------------------------------
// utility functions
// ------------------------------------------------------------------------------------------------

func getClassUtilsFromFile(srcdir, currentPath, boFileName string) (clsuCore *classUtilsCore) {
	// controlling the file
	filename := path.Join(srcdir, currentPath, boFileName)
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
						if clsuCore == nil {
							clsuCore = &classUtilsCore{
								class:     className(typeSpec.Name.Name),
								lastBOMod: stat.ModTime(),
								srcPath:   currentPath,
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
	currentModuleName moduleName
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
func getCurrentModuleName() moduleName {
	if currentModuleName == "" {
		currentModuleName = moduleName(path.Base(getCurrentModule()))
	}

	return currentModuleName
}

// ------------------------------------------------------------------------------------------------
// generating the ClassUtils (*--clsu.go) files
// ------------------------------------------------------------------------------------------------

const clsuFileTemplate = `// Generated file, do not edit!
package %s

import (
	"github.com/aldesgroup/goald"
	"%s"
)

type $$CLASSNAME$$ClassUtils struct {
	goald.IClassUtilsCore
}

func ClassUtilsFor$$CLASSNAME$$(srcPath, lastMod string) goald.IClassUtils {
	return &$$CLASSNAME$$ClassUtils{IClassUtilsCore: goald.NewClassUtilsCore(srcPath, lastMod)}
}

func (thisUtils *$$CLASSNAME$$ClassUtils) NewObject() any {
	return &$$PKG$$.$$CLASSNAME$${}
}

func (thisUtils *$$CLASSNAME$$ClassUtils) NewSlice() any {
	return []*$$PKG$$.$$CLASSNAME$${}
}
`

func genClassUtilsFile(srcdir string, clsuCore *classUtilsCore) {
	// the class utils filename
	clsuFilename := path.Join(srcdir, clsuCore.srcPath, sourceCLASSxUTILSxDIR,
		fmt.Sprintf("%s%s", utils.PascalToKebab(string(clsuCore.class)), sourceCLSUxSUFFIX))

	// does it exist?
	if !utils.FileExists(clsuFilename) {
		slog.Info(fmt.Sprintf("Will generate utils: %s", clsuFilename))
		importPkg := path.Join(getCurrentModule(), clsuCore.srcPath)
		content := fmt.Sprintf(clsuFileTemplate, sourceCLASSxUTILSxDIR, importPkg)
		content = strings.ReplaceAll(content, "$$CLASSNAME$$", string(clsuCore.class))
		content = strings.ReplaceAll(content, "$$PKG$$", path.Base(importPkg))
		utils.WriteToFile(content, clsuFilename)
	}
}
