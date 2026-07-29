// ------------------------------------------------------------------------------------------------
// Here is the code used to generate the classes, which is the static information we can extract
// from using AST/basic code parsing.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const includePATH = "_include"
const sourceFILExSUFFIX = "--.go"
const sourceCLSxSUFFIX = "--cls.go"
const sourceREGISTRYxNAME = "registry.go"
const sourceCLASSxDIR = "class"

// ------------------------------------------------------------------------------------------------
// Going over all the physical source code files and generating stuff along the way
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateAllClasses(srcdir, currentPath string, _ bool,
	allEntriesInCodeSoFar map[packageName]map[utils.ClassName]*baseClass, regen bool) (codeChanged bool) {
	// we want all the entities we find in the code to build 1 global registry
	allClassesInCode := allEntriesInCodeSoFar

	// the path we're currently reading at e.g. go/pkg1/pkg2
	readingPath := path.Join(srcdir, currentPath)

	// are we currently dealing with a package with business objects ?
	var currentPackage packageName

	// going through the resources found withing the current directory
	for _, entry := range core.EnsureReadDir(readingPath) {
		if entry.IsDir() {
			// not going into the vendor, nor the git folder obviously
			if entry.Name() != "vendor" && entry.Name() != ".git" {
				// found another directory, let's dive deeper!
				codeChanged = thisServer.generateAllClasses(srcdir, path.Join(currentPath, entry.Name()), false, allClassesInCode, regen) || codeChanged
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with sourceFILExSUFFIX
			if strings.HasSuffix(entry.Name(), sourceFILExSUFFIX) {
				// hadn't we figured out yet we're dealing with a BO package?
				if currentPackage == "" {
					// now we have
					currentPackage = packageName(path.Base(currentPath))

					// but at this point the package should not exist yet, or it means we have 2 packages with the same name
					core.PanicMsgIf(allEntriesInCodeSoFar[currentPackage] != nil, "there are 2 packages named %s which is not allowed!", currentPackage)
					allEntriesInCodeSoFar[currentPackage] = map[utils.ClassName]*baseClass{}
					thisServer.Warn("Found new package " + string(currentPackage))
				}

				// getting the business object entry for the egustry, from the current file
				if class := getClassFromFile(srcdir, currentPath, entry.Name()); class != nil {
					// checking the biz obj / file naming
					if expected := core.PascalToKebab(string(class.class)) + sourceFILExSUFFIX; expected != entry.Name() {
						core.PanicMsg("The business object's name should be the file name Pascal-cased, i.e. we should have: "+
							"%s in file %s, "+
							"or %s in file %s",
							class.class, expected,
							core.KebabToPascal(strings.Replace(entry.Name(), sourceFILExSUFFIX, "", 1)), entry.Name(),
						)
					}

					// checking the unicity of each biz obj name
					if allClassesInCode[currentPackage][class.class] != nil {
						core.PanicMsg("We can't have 2 business objects with the same name '%s'."+
							" This would lead to the same REST path. You have to rename one.", class.class)
					} else {
						// adding one more BO to our list
						allClassesInCode[currentPackage][class.class] = class

						// generating the corresponding Class file, if it doesn't exist yet
						codeChanged = thisServer.genClassFile(srcdir, class, regen) || codeChanged
					}
				} else {
					thisServer.Error(false, "No business object found in file "+entry.Name())
				}
			}
		}
	}

	// if we're at root here, this means we've browsed through all the code already,
	// and can now decide to (re-)generate the object registry - or not
	if currentPath == "." {
		codeChanged = thisServer.writeRegistryFilesIfNeeded(srcdir, allClassesInCode, regen) || codeChanged
	}

	return
}

// ------------------------------------------------------------------------------------------------
// Writing the registry file
// ------------------------------------------------------------------------------------------------

const goaldIMPORT = "g \"github.com/aldesgroup/goald\""
const registryFileTemplate = `// Generated file, do not edit!
package %s

import (
	` + goaldIMPORT + `
%s
)

func init() {
%s
}
`

func (thisServer *server) writeRegistryFilesIfNeeded(srcdir string, allClassesInCode map[packageName]map[utils.ClassName]*baseClass, regen bool) (codeChanged bool) {
	// do we need to regenerate the object registry at the current path?
	needRegen := regen

	// iterating over all the packages we've found
	for currentPackage, allClassesInPackage := range allClassesInCode {
		// let's check the current class, the ones coded right now
		for clsName, classInCode := range allClassesInPackage {
			classInRegistry := classFor(clsName, false)
			if classInRegistry == nil {
				thisServer.Info(fmt.Sprintf("Business object '%s' has appeared since the last generation!", clsName))
				needRegen = true

				break
			} else if classInRegistry.getLastBOMod().Before(classInCode.getLastBOMod()) {
				thisServer.Info(fmt.Sprintf("Business object '%s' has changed since the last generation!", clsName))
				needRegen = true

				break
			}
		}

		// if we're not doing regen because of added or changed biz objs,
		// maybe we have to because of deleted ones!
		if !needRegen {
			for clsName, class := range classRegistry.items {
				if class.getModule() == getCurrentModuleName() &&
					class.getPackage() == string(currentPackage) &&
					allClassesInPackage[clsName] == nil {

					needRegen = true

					thisServer.Info(fmt.Sprintf("Business object '%s' has disappeared since the last generation!", clsName))

					break
				}
			}
		}

		// now let's write the registry file, if needed, and if we're at root
		if nbEntries := len(allClassesInPackage); nbEntries > 0 && needRegen {
			// gathering the biz objs in order
			registrationLines := []string{fmt.Sprintf("\tg.In(\"%s\")", getCurrentModuleName())}

			// and the imports, but only once per import, hence the map
			imports := []string{}
			imported := map[string]bool{}

			// going over all the class cores
			for _, class := range core.GetSortedValues(allClassesInPackage) {
				// adding 1 registration line per business object
				boPath := path.Base(class.srcPath)
				registrationLines = append(registrationLines,
					fmt.Sprintf("%sRegister(%s.ClassFor%s(\"%s\", \"%s\"))", "\t\t", boPath, class.class,
						class.srcPath, class.getLastBOMod().Add(time.Second).Format(time.RFC3339)),
				)

				// adding the corresponding import
				if !imported[class.srcPath] {
					imports = append(imports, boPath+" \""+path.Join(getCurrentModule(), includePATH, boPath, sourceCLASSxDIR)+"\"")
					imported[class.srcPath] = true
				}
			}

			// which package?
			pkgName := string(currentPackage)

			// which file?
			filename := path.Join(srcdir, includePATH, pkgName, sourceREGISTRYxNAME)

			// which content?
			dot := "." + newline
			content := fmt.Sprintf(registryFileTemplate, pkgName, strings.Join(imports, newline), strings.Join(registrationLines, dot))

			// writing to the file
			core.WriteToFile(content, filename)

			codeChanged = true
		}
	}

	return
}

// ------------------------------------------------------------------------------------------------
// utility functions
// ------------------------------------------------------------------------------------------------

func getClassFromFile(srcdir, currentPath, boFileName string) (class *baseClass) {
	// controlling the file
	filename := path.Join(srcdir, currentPath, boFileName)
	stat, errStat := os.Stat(filename)
	core.PanicMsgIfErr(errStat, "Could not check the modification time for file '%s'", filename)

	// parsing the file to get the AST (Abstract Syntax Tree)
	file, errParse := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	core.PanicMsgIfErr(errParse, "Error while parsing '%s'", filename)

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
					switch objType := typeSpec.Type.(type) {
					case *ast.StructType, *ast.InterfaceType:
						if class == nil {
							_, isInterface := objType.(*ast.InterfaceType)
							class = &baseClass{
								class:     utils.ClassName(typeSpec.Name.Name),
								lastBOMod: stat.ModTime(),
								srcPath:   currentPath,
								intrface:  isInterface,
							}
						} else {
							core.PanicMsg("More than one struct declared in the BusinessObject file '%s'!", filename)
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
			core.PanicMsg("Could not read the build infos!")
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
// generating the Class (*--class.go) files
// ------------------------------------------------------------------------------------------------

const classFileTemplateBase = `// Generated file, do not edit!
package %s

import (
	"github.com/aldesgroup/goald"%s
)

type $$CLASSNAME$$Class struct {
	goald.IClass
}

func ClassFor$$CLASSNAME$$(srcPath, lastMod string) goald.IClass {
	return &$$CLASSNAME$$Class{IClass: goald.NewClass(srcPath, "$$CLASSNAME$$", lastMod)%s}
}
`

const classFileTemplateConcrete = `
func (thisClass *$$CLASSNAME$$Class) NewObject() any {
	return &$$PKG$$.$$CLASSNAME$${}
}

func (thisClass *$$CLASSNAME$$Class) NewSlice() any {
	return []*$$PKG$$.$$CLASSNAME$${}
}
`

const classFileTemplateInterface = `
func (thisClass *$$CLASSNAME$$Class) NewObject() any {
	panic("NewObject cannot be called for an interface!")
}

func (thisClass *$$CLASSNAME$$Class) NewSlice() any {
	panic("NewSlice cannot be called for an interface!")
}
`

func (thisServer *server) genClassFile(srcdir string, class *baseClass, regen bool) (codeChanged bool) {
	// the class filename - the "class" folder lives as a sibling of the "model" folder,
	// inside the feature package's folder in the "_include" directory
	classFilename := path.Join(srcdir, includePATH, path.Base(class.srcPath), sourceCLASSxDIR,
		core.PascalToKebab(string(class.class))+sourceCLSxSUFFIX)

	// does it exist?
	if !core.FileExists(classFilename) || regen {
		thisServer.Info(fmt.Sprintf("Will generate class: %s", classFilename))
		content := classFileTemplateBase
		importPkg := path.Join(getCurrentModule(), class.srcPath)
		toImport := ""
		asInterface := ""
		if class.isInterface() {
			content += classFileTemplateInterface
			asInterface = ".AsInterface()"
		} else {
			content += classFileTemplateConcrete
			toImport = fmt.Sprintf("\n\"%s\"", importPkg)
		}
		content = fmt.Sprintf(content, sourceCLASSxDIR, toImport, asInterface)
		content = strings.ReplaceAll(content, "$$CLASSNAME$$", string(class.class))
		content = strings.ReplaceAll(content, "$$PKG$$", path.Base(importPkg))
		core.WriteToFile(content, classFilename)
		return true
	}

	return false
}
