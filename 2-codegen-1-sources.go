// ------------------------------------------------------------------------------------------------
// Here is the code used to generate the BO sources, which is the static information we can extract
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
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

const includePATH = "_include"
const boFILExSUFFIX = "--.go"
const sourceFILExSUFFIX = "--src.go"
const registryFILExNAME = "registry.go"
const sourceFOLDERxNAME = "source"

// ------------------------------------------------------------------------------------------------
// Going over all the physical source code files and generating stuff along the way
// ------------------------------------------------------------------------------------------------

func (thisServer *server) generateAllSources(srcdir, currentPath string, _ bool, regen bool,
	allSourcesInCodeSoFar map[packageName]map[utils.ModelName]*baseBusinessObjectModelSource) (codeChanged bool) {
	// we want all the entities we find in the code to build 1 global registry
	allSourcesInCode := allSourcesInCodeSoFar

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
				codeChanged = thisServer.generateAllSources(srcdir, path.Join(currentPath, entry.Name()), false, regen, allSourcesInCode) || codeChanged
			}
		} else {
			// found a file... but we're only interested in files containing Business Objects, which must end with boFILExSUFFIX
			if strings.HasSuffix(entry.Name(), boFILExSUFFIX) {
				// hadn't we figured out yet we're dealing with a BO package?
				if currentPackage == "" {
					// now we have
					currentPackage = packageName(path.Base(currentPath))

					// but at this point the package should not exist yet, or it means we have 2 packages with the same name
					core.PanicMsgIf(allSourcesInCodeSoFar[currentPackage] != nil, "there are 2 packages named %s which is not allowed!", currentPackage)
					allSourcesInCodeSoFar[currentPackage] = map[utils.ModelName]*baseBusinessObjectModelSource{}
					thisServer.Warn("Found new package " + string(currentPackage))
				}

				// getting the business object entry for the egustry, from the current file
				if source := thisServer.getSourceFromFile(srcdir, currentPath, entry.Name()); source != nil {
					// checking the biz obj / file naming
					if expected := core.PascalToKebab(string(source.modelName)) + boFILExSUFFIX; expected != entry.Name() {
						core.PanicMsg("The business object's name should be the file name Pascal-cased, i.e. we should have: "+
							"%s in file %s, "+
							"or %s in file %s",
							source.modelName, expected,
							core.KebabToPascal(strings.Replace(entry.Name(), boFILExSUFFIX, "", 1)), entry.Name(),
						)
					}

					// checking the unicity of each biz obj name
					if allSourcesInCodeSoFar[currentPackage][source.modelName] != nil {
						core.PanicMsg("We can't have 2 business objects with the same name '%s'."+
							" This would lead to the same REST path. You have to rename one.", source.modelName)
					} else {
						// adding one more BO to our list
						allSourcesInCodeSoFar[currentPackage][source.modelName] = source

						// generating the corresponding source file, if it doesn't exist yet
						codeChanged = thisServer.genSourceFile(srcdir, source, regen) || codeChanged
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
		codeChanged = thisServer.writeRegistryFilesIfNeeded(srcdir, regen, allSourcesInCodeSoFar) || codeChanged
	}

	return
}

// ------------------------------------------------------------------------------------------------
// utility functions
// ------------------------------------------------------------------------------------------------

func (thisServer *server) getSourceFromFile(srcdir, currentPath, boFileName string) (source *baseBusinessObjectModelSource) {
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
						if source == nil {
							_, isInterface := objType.(*ast.InterfaceType)
							source = &baseBusinessObjectModelSource{
								modelName: utils.ModelName(typeSpec.Name.Name),
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

func (thisServer *server) getModelFromFile(srcdir, currentPath, boFileName string) IBusinessObjectModel {
	return modelFor(thisServer.getSourceFromFile(srcdir, currentPath, boFileName).getName(), true)
}

// ------------------------------------------------------------------------------------------------
// generating the source (*--src.go) files
// ------------------------------------------------------------------------------------------------

const sourceFileTemplateBase = `// Generated file, do not edit!
package %[1]s

import (
	"github.com/aldesgroup/goald"%[3]s
)

type %[5]sModelSource struct {
	goald.IBusinessObjectModelSource
}

func For%[5]s(srcPath, lastMod string) goald.IBusinessObjectModelSource {
	return &%[5]sModelSource{IBusinessObjectModelSource: goald.NewBusinessObjectModelSource(srcPath, "%[5]s", lastMod)%[4]s}
}
`

const sourceFileTemplateConcrete = `
func (this *%[5]sModelSource) NewObject() any {
	return &%[6]s.%[5]s{}
}

func (this *%[5]sModelSource) NewSlice() any {
	return []*%[6]s.%[5]s{}
}
`

const sourceFileTemplateInterface = `
func (this *%[5]sModelSource) NewObject() any {
	panic("NewObject cannot be called for an interface!")
}

func (this *%[5]sModelSource) NewSlice() any {
	panic("NewSlice cannot be called for an interface!")
}
`

func (thisServer *server) genSourceFile(srcdir string, source *baseBusinessObjectModelSource, regen bool) (codeChanged bool) {
	// the source filename - the "source" folder lives as a sibling of the "model" folder,
	// inside the feature package's folder in the "_include" directory
	sourceFilename := path.Join(srcdir, includePATH, path.Base(source.srcPath), sourceFOLDERxNAME,
		core.PascalToKebab(string(source.modelName))+sourceFILExSUFFIX)

	// does it exist?
	if !core.FileExists(sourceFilename) || regen {
		thisServer.Info(fmt.Sprintf("Will generate source: %s", sourceFilename))
		content := sourceFileTemplateBase
		importPkg := path.Join(getCurrentSourceModule(), source.srcPath)
		toImport := ""
		asInterface := ""
		if source.isInterface() {
			content += sourceFileTemplateInterface
			asInterface = ".AsInterface()"
		} else {
			content += sourceFileTemplateConcrete
			toImport = fmt.Sprintf("\n\"%s\"", importPkg)
		}
		content = fmt.Sprintf(content,
			sourceFOLDERxNAME,        // 1
			getCurrentSourceModule(), // 2
			toImport,                 // 3
			asInterface,              // 4
			source.getName(),         // 5
			path.Base(importPkg),     // 6
		)
		core.WriteToFile(content, sourceFilename)
		return true
	}

	return false
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

func (thisServer *server) writeRegistryFilesIfNeeded(srcdir string, regen bool,
	allSourcesInCode map[packageName]map[utils.ModelName]*baseBusinessObjectModelSource) (codeChanged bool) {
	// do we need to regenerate the object registry at the current path?
	needRegen := regen

	// the needed DB connections
	neededDBs := map[string]bool{}

	// iterating over all the packages we've found
	for currentPackage, allSourcesInPackage := range allSourcesInCode {
		// let's check the current sources in code, the ones coded right now
		for modelName, sourceInCode := range allSourcesInPackage {
			sourceInRegistry := sourceRegistry.items[modelName]
			if sourceInRegistry == nil {
				thisServer.Info(fmt.Sprintf("Business object '%s' has appeared since the last generation!", modelName))
				needRegen = true
			} else if sourceInRegistry.getLastBOMod().Before(sourceInCode.getLastBOMod()) {
				thisServer.Info(fmt.Sprintf("Business object '%s' has changed since the last generation!", modelName))
				needRegen = true
			}

			// taking the opportunity to check for the needed DB connections, if there's any info about it yet
			if model := modelFor(modelName, false); model != nil {
				if db := model.getDB(); db != nil && db.schema != nil {
					neededDBs[string(db.schema.DbServer.Type)] = true
				}
			}
		}

		// if we're not doing regen because of added or changed biz objs,
		// maybe we have to because of deleted ones!
		if !needRegen {
			for modelName, source := range sourceRegistry.items {
				if source.getModule() == getCurrentSourceModuleName() &&
					source.getPackage() == string(currentPackage) &&
					allSourcesInPackage[modelName] == nil {

					needRegen = true

					thisServer.Info(fmt.Sprintf("Business object '%s' has disappeared since the last generation!", modelName))

					break
				}
			}
		}

		// now let's write the registry file, if needed, and if we're at root
		if nbEntries := len(allSourcesInPackage); nbEntries > 0 && needRegen {
			// gathering the biz objs in order
			registrationLines := []string{fmt.Sprintf("\tg.In(\"%s\")", getCurrentSourceModuleName())}

			// and the imports, but only once per import, hence the map
			imports := []string{}
			imported := map[string]bool{}

			// going over all the sources
			for _, source := range core.GetSortedValues(allSourcesInPackage) {
				// adding 1 registration line per business object
				boPath := path.Base(source.srcPath)
				registrationLines = append(registrationLines,
					fmt.Sprintf("%sRegister(source.For%s(\"%s\", \"%s\"))", "\t\t", source.modelName,
						source.srcPath, source.getLastBOMod().Add(time.Second).Format(time.RFC3339)),
				)

				// adding the corresponding import
				if !imported[source.srcPath] {
					imports = append(imports, "\""+path.Join(string(getCurrentSourceModule()), includePATH, boPath, sourceFOLDERxNAME)+"\"")
					imports = append(imports, "_ \""+path.Join(string(getCurrentSourceModule()), includePATH, boPath, modelFOLDERxNAME)+"\"")
					for dbFolder := range neededDBs {
						imports = append(imports, "_ \""+path.Join(string(getCurrentSourceModule()), includePATH, dbFOLDERNAME, dbFolder)+"\"")
					}
					imported[source.srcPath] = true
				}
			}

			// which file?
			filename := path.Join(srcdir, includePATH, string(currentPackage), registryFILExNAME)

			// which content?
			dot := "." + newline
			content := fmt.Sprintf(registryFileTemplate, string(currentPackage), strings.Join(imports, newline), strings.Join(registrationLines, dot))

			// writing to the file
			core.WriteToFile(content, filename)

			// since we're importing the /model package, we make sure this package can be imported, i.e. it exists and has an index.go file
			core.WriteStringToFile(path.Join(srcdir, includePATH, string(currentPackage), modelFOLDERxNAME, "index.go"),
				"// Generated file, do not edit!"+newline+"package model")

			codeChanged = true
		}
	}

	return
}
