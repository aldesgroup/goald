// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the BO models in the web app
// ------------------------------------------------------------------------------------------------
package goald

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// constants, variables, useful structs... & main generation function
// ------------------------------------------------------------------------------------------------

var (
	modelsDIRPATH  = path.Join("src", "components", "models")
	handledClasses = map[className]bool{}
)

const (
	newFieldNAME = "newField"
	newModelNAME = "newModel"
)

// TODO maybe do not collocate everything on the server... let's change the receiver here

func (thisServer *server) generateAllClientAppModels(destdir string, regen bool, isWebapp bool) {
	// the enum files to generate
	enums := map[string]IEnum{}

	// scanning for BOs involved in endpoints used from the web app
	for _, ep := range restRegistry.endpoints {
		if (isWebapp && ep.isCalledFromWebApp()) || (!isWebapp && ep.isCalledFromNativeApp()) {
			// generating the model for the endpoint resource
			thisServer.generateClientAppModel(destdir, ep, false, enums, regen, isWebapp)

			// if the endpoint admits a BO as an input (body or URL params), then we also need the model in the webapp
			if ep.getInputOrParamsClass() != "" {
				thisServer.generateClientAppModel(destdir, ep, true, enums, regen, isWebapp)
			}
		}
	}

	// enum files generation
	for enumType, enum := range enums {
		// TODO run in go routines
		generateWebAppEnum(destdir, enumType, enum)
	}
}

type codeContext struct {
	enums      map[string]IEnum
	bObjType   utils.GoaldType
	boInstance utils.GoaldValue
}

func (ctx *codeContext) getEnumType(field IField) string {
	return ctx.bObjType.FieldByName(field.getName()).Type().Name()
}

func (thisServer *server) generateClientAppModel(destdir string, ep iEndpoint, useInputClass bool,
	enums map[string]IEnum, regen bool, isWebapp bool) {
	// which model to generate?
	clsName := core.IfThenElse(useInputClass, ep.getInputOrParamsClass(), ep.getResourceClass())

	// already done this?
	if handledClasses[clsName] {
		return
	}

	// but we're doing this now
	handledClasses[clsName] = true

	// the business object we're dealing with
	boSpecs := specsForName(clsName)
	boClass := getClass(boSpecs)
	boFields := core.GetSortedValues(boSpecs.base().fields)

	// the file we're dealing with
	modelName := core.PascalToCamel(string(clsName))
	filename := modelName + ".ts"
	filepath := path.Join(destdir, modelsDIRPATH, filename)

	// gathering needed info into a context
	codeCtx := &codeContext{
		enums:      enums,
		bObjType:   utils.TypeOf(boClass.NewObject(), true),
		boInstance: utils.ValueOf(boClass.NewObject()),
	}

	// gathering the needed enums
	for _, field := range boFields {
		if field.getTypeFamily() == utils.TypeFamilyENUM {
			enums[codeCtx.getEnumType(field)] = codeCtx.boInstance.GetFieldValue(field.getName()).(IEnum)
		}
	}

	// do we need to (re)generate the file?
	if !regen && core.FileExists(filepath) && core.EnsureModTime(filepath).After(boClass.getLastBOMod()) {
		return // the file already exists and is older than our changes in the BO class file
	}

	// getting the file content - which might be empty if the file does not exist yet
	code := parseCode(filepath).initFixedBlocks(modelName, thisServer.config.commonPart().HTTP.ApiPath+ep.getFullPath(), isWebapp)

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range boFields {
		code.addFieldIfNeeded(codeCtx, field)
	}

	// "unpacking" the code blocks to code lines
	codeLines := []string{}
	for _, block := range code.blocks {
		// the block should at least have a non-nil code line
		if block.lines[0] != nil {
			for _, line := range block.lines {
				if line != nil {
					codeLines = append(codeLines, line.rawline)
				}
			}
		}
	}

	// writing out the code lines
	core.WriteToFile(strings.Join(codeLines, newline), filepath)
	slog.Info(fmt.Sprintf("(Re-)generated file %s", filepath))
}

// ------------------------------------------------------------------------------------------------
// code editing
// ------------------------------------------------------------------------------------------------

// initiating the code, with its structure, i.e. with its required imports, and the fixed variables / functions
func (thisCode *codeFile) initFixedBlocks(modelName string, endpointPath string, isWebapp bool) *codeFile {
	if len(thisCode.blocks) == 0 {
		if isWebapp {
			thisCode.addNewBlock("import { newField, newModel } from '~/vendor/goaldr';", true, "", false)
		} else {
			thisCode.addNewBlock("import { newField, newModel } from '~/vendor/goaldn';", true, "", false)
		}
		thisCode.addNewBlock("export const "+modelName+" = "+newModelNAME+"('"+endpointPath+"', {", true, newModelNAME, true).appendLine("});", true)
	}

	return thisCode
}

// handling a field, adding it if not in the code already, flagging an enum for generation if it's an enum field
func (thisCode *codeFile) addFieldIfNeeded(codeCtx *codeContext, field IField) {
	// adding to the context, and the class file content
	if typeFamily := field.getTypeFamily(); typeFamily != utils.TypeFamilyUNKNOWN && typeFamily != utils.TypeFamilyRELATIONSHIPxMONOM {
		// not handling multiple properties for now - nor the ID field
		if !field.isMultiple() && field.getName() != "ID" {
			var (
				enumType, enumVar, initVal, fieldAtomType string
			)

			// dealing with some field specificities
			switch typeFamily {
			// --- enums -------------------------------------------------------------------
			case utils.TypeFamilyENUM:
				// flagging this enum type for code generation
				enumType = codeCtx.getEnumType(field)

				// importing the enum if needed
				enumVar = "_" + core.PascalToCamel(enumType)
				if thisCode.blocksMap[enumVar] == nil {
					// panic(fmt.Sprintf("Missing import for '%s'; should add it after block %d", enumType, thisCode.findLastImportPosition()))
					thisCode.addEnumImport(enumVar)
				}

				// setting the field atom's type
				fieldAtomType = fmt.Sprintf("<%s.%s>", enumVar, enumType)

				// proposing an init value
				initVal = fmt.Sprintf("%s.%s", enumVar, makeEnumName(core.GetFirstMapValue(codeCtx.enums[enumType].Values())))

			// --- numbers -----------------------------------------------------------------
			case utils.TypeFamilyINT, utils.TypeFamilyBIGINT, utils.TypeFamilyREAL, utils.TypeFamilyDOUBLE:
				// proposing an init value
				initVal = "0"

			// --- booleans ----------------------------------------------------------------
			case utils.TypeFamilyBOOL:
				// proposing an init value
				initVal = "false"

			// --- dates ----------------------------------------------------------------
			case utils.TypeFamilyDATE:
				// setting the field atom's type
				fieldAtomType = "<Date | null>"

				// proposing an init value
				initVal = "null"

				// SWITCH END
			}

			// adding the field name to the model block if needed
			if !thisCode.blockHasLineStartingWith(newModelNAME, field.getName()+":") {
				thisCode.insertLineIntoBlockBeforePrefix(newModelNAME, fmt.Sprintf("    %s,", field.getName()), "}")
			}

			// adding the field if needed
			missingField := thisCode.blocksMap[field.getName()] == nil
			if missingField {
				fieldDecl := fmt.Sprintf("const %s = "+newFieldNAME+"%s('%s', {", field.getName(), fieldAtomType, field.getName())
				newBlock := thisCode.addNewBlockBeforeEndPosition(fieldDecl, true, field.getName(), true, 1)
				newBlock.appendLine(fmt.Sprintf("    initialValue: %s?,", initVal), true)
				newBlock.appendLine("});", true)
			}

			// linking the enum's options to the field, if needed
			if typeFamily == utils.TypeFamilyENUM {
				if missingField || !thisCode.blockHasLineStartingWith(field.getName(), "options:") {
					thisCode.insertLineIntoBlockBeforePrefix(field.getName(), fmt.Sprintf("    options: %s.Options,", enumVar), "}")
				}

				if enumField := field.(*EnumField); len(enumField.onlyValues) > 0 {
					restrictedValues := []string{}
					for _, restrictedValue := range enumField.onlyValues {
						restrictedValues = append(restrictedValues, fmt.Sprintf("%s.%s", enumVar, makeEnumName(restrictedValue.String())))
					}
					thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    optionsOnly: [%s],", strings.Join(restrictedValues, ", ")), "optionsOnly:", "}")
				}
			}

			// handling the constraints - for numeric fields
			if numField, ok := field.(iNumericField); ok {
				// min constraint
				if numField.isMinSet() {
					switch nf := numField.(type) {
					case *IntField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    min: %d,", nf.min), "min:", "}")
					case *BigIntField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    min: %d,", nf.min), "min:", "}")
					case *RealField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    min: %f,", nf.min), "min:", "}")
					case *DoubleField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    min: %f,", nf.min), "min:", "}")
					}
				}

				// max constraint
				if numField.isMaxSet() {
					switch nf := numField.(type) {
					case *IntField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    max: %d,", nf.max), "max:", "}")
					case *BigIntField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    max: %d,", nf.max), "max:", "}")
					case *RealField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    max: %f,", nf.max), "max:", "}")
					case *DoubleField:
						thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    max: %f,", nf.max), "max:", "}")
					}
				}
			}

			// handling the constraints - for string fields
			if sf, ok := field.(*StringField); ok {
				if sf.size > 0 {
					thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    max: %d,", sf.size), "max:", "}")
				}
				if sf.atLeast > 0 {
					thisCode.updateLineIntoBlockWithPrefix(field.getName(), fmt.Sprintf("    min: %d,", sf.atLeast), "min:", "}")
				}
			}

			// handling the constraints - misc
			if field.isMandatory() {
				thisCode.updateLineIntoBlockWithPrefix(field.getName(), "    mandatory: true,", "mandatory:", "}")
			} else {
				thisCode.updateLineIntoBlockWithPrefix(field.getName(), "    mandatory: false,", "mandatory: true", "")
			}
		}
	}
}

// adding an import for an enum
func (thisCode *codeFile) addEnumImport(enumVar string) {
	enumImport := fmt.Sprintf("import * as %s from \"./%s\"", enumVar, enumVar)
	thisCode.addNewBlockAtPosition(enumImport, true, enumVar, false, thisCode.findLastImportPosition())
}

func (thisCode *codeFile) findLastImportPosition() (pos int) {
	for pos = 1; pos < len(thisCode.blocks); pos++ {
		if !strings.HasPrefix(thisCode.blocks[pos].lines[0].rawline, "import ") {
			break
		}
	}

	return
}

// ------------------------------------------------------------------------------------------------
// enum files generation
// ------------------------------------------------------------------------------------------------

func generateWebAppEnum(destdir string, enumType string, enum IEnum) {
	filename := "_" + core.PascalToCamel(enumType) + ".ts"
	filepath := path.Join(destdir, modelsDIRPATH, filename)

	content := "// Generated by Aldev, do not edit!" + newline
	allTypes := []string{}
	typeDecl := []string{}
	for _, enumVal := range core.GetSortedKeys(enum.Values()) {
		enumLabel := enum.Values()[enumVal]
		enumName := makeEnumName(enumLabel)
		content += fmt.Sprintf("export const %s = %d;", enumName, enumVal) + newline
		allTypes = append(allTypes, fmt.Sprintf("    { value: %s, label: \"%s\" },", enumName, enumLabel))
		typeDecl = append(typeDecl, fmt.Sprintf("typeof %s", enumName))
	}

	content += "export const Options = [" + newline
	content += strings.Join(allTypes, newline) + newline
	content += "];" + newline
	content += fmt.Sprintf("export type %s = %s;", enumType, strings.Join(typeDecl, " | "))

	core.WriteToFile(content, filepath)
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

// ------------------------------------------------------------------------------------------------
// code parsing & generic code editing
// ------------------------------------------------------------------------------------------------

type codeLine struct {
	rawline string
	isCode  bool
	content string
}

func (line *codeLine) with(content string) *codeLine {
	line.content = content
	return line
}

// func newCodeLine(content string, num int, isCode bool) *codeLine {
func newCodeLine(rawline string, isCode bool) *codeLine {
	return &codeLine{
		rawline: rawline,
		isCode:  isCode,
	}
}

type codeBlock struct {
	// we should be able to insert new lines, for new min / max constraints for instance, or new field atoms in the form
	lines []*codeLine
	id    string
}

// how to add a line to a code block
func (block *codeBlock) appendLine(rawline string, isCode bool) *codeBlock {
	block.lines = append(block.lines, newCodeLine(rawline, isCode))
	return block
}

type codeFile struct {
	blocks    []*codeBlock
	blocksMap map[string]*codeBlock // some code blocks mapped with their ID
	current   *codeBlock            // used during the parsing
}

// adding a new block of code after the last one already present
func (thisCode *codeFile) addNewBlock(rawline string, isCode bool, id string, blankBefore bool) *codeBlock {
	return thisCode.addNewBlockAtPosition(rawline, isCode, id, blankBefore, -1)
}

// adding a block after another one
func (thisCode *codeFile) addNewBlockAfterBlock(rawline string, isCode bool, id string, blankBefore bool, blockID string) *codeBlock {
	// finding the targeted block
	blockPos := 0
	for pos, block := range thisCode.blocks {
		if block.id == blockID {
			blockPos = pos + 1
			break
		}
	}

	return thisCode.addNewBlockAtPosition(rawline, isCode, id, blankBefore, blockPos)
}

// adding a new block of code before the one counted from the end, with indexFromEnd, which must be >= 1;
// if we use indexFromEnd = 1, then we're inserting something before the last code block;
// if we use indexFromEnd = 2, then it's before the before-the-last; etc...
func (thisCode *codeFile) addNewBlockBeforeEndPosition(rawline string, isCode bool, id string, blankBefore bool, indexFromEnd int) *codeBlock {
	return thisCode.addNewBlockAtPosition(rawline, isCode, id, blankBefore, len(thisCode.blocks)-indexFromEnd)
}

// adding a new code block starting with the given line, possibly with an ID; if the given position is > 0,
// then the block is not appended to the others, but inserted between them at the desired position
func (thisCode *codeFile) addNewBlockAtPosition(rawline string, isCode bool, id string, blankBefore bool, pos int) *codeBlock {
	// if there was a block before the one we're starting, we want to get the lines that follows it
	// - like comments - to attach them to the new block
	newBlockNonCodeLines := []*codeLine{}
	if thisCode.current != nil {
		// checking each line of the current block from the end
		for i := len(thisCode.current.lines) - 1; i >= 0; i-- {
			if line := thisCode.current.lines[i]; !line.isCode {
				// if it's non code, it shall be moved from the current block...
				thisCode.current.lines[i] = nil
				// ... to the next block
				newBlockNonCodeLines = append([]*codeLine{line}, newBlockNonCodeLines...)
			} else {
				break
			}
		}
	}

	// new current block
	thisCode.current = &codeBlock{
		lines: newBlockNonCodeLines,
		id:    id,
	}

	// adding a blank line if needed
	if blankBefore {
		thisCode.current.appendLine("", isCode)
	}

	// adding the code line
	thisCode.current.appendLine(rawline, isCode)

	// adding the block
	if pos < 0 {
		// no desired position = just append aftr the other
		thisCode.blocks = append(thisCode.blocks, thisCode.current)
	} else {
		thisCode.blocks = append(thisCode.blocks[:pos], append([]*codeBlock{thisCode.current}, thisCode.blocks[pos:]...)...)
	}

	// mapping the block with it's ID, if one is provided
	if id != "" {
		thisCode.blocksMap[id] = thisCode.current
	}

	// returning the new block
	return thisCode.current
}

// appending a line to the block targeted by the given ID, before the line that starts with the given prefix
func (thisCode *codeFile) insertLineIntoBlockBeforePrefix(blockID string, newLine string, insertBeforePrefix string) {
	if block := thisCode.blocksMap[blockID]; block != nil {
		// finding where to cut the block lines
		var pos int
		for pos = len(block.lines) - 1; pos >= 0; pos-- {
			line := block.lines[pos]
			if line != nil && strings.HasPrefix(line.rawline, insertBeforePrefix) {
				break
			}
		}

		// inserting the new line at the right position
		if pos >= 0 {
			block.lines = append(block.lines[:pos], append([]*codeLine{newCodeLine(newLine, true)}, block.lines[pos:]...)...)
		} else {
			slog.Error(fmt.Sprintf("Could not insert '%s' into block '%s' before '%s'. Block =", newLine, blockID, insertBeforePrefix))
			for _, line := range block.lines {
				slog.Error(line.rawline)
			}
		}
	} else {
		slog.Error("No block found with ID: " + blockID)
	}
}

// // inserting a line to the block targeted by the given ID, immediately after the line that starts with the given prefix
// func (thisCode *codeFile) insertLineIntoBlockAfterPrefix(blockID string, newLine string, insertAfterPrefix string) {
// 	if block := thisCode.blocksMap[blockID]; block != nil {
// 		for pos := 0; pos < len(block.lines); pos++ {
// 			line := block.lines[pos]
// 			if line != nil && strings.HasPrefix(line.rawline, insertAfterPrefix) {
// 				insertPos := pos + 1
// 				block.lines = append(block.lines[:insertPos], append([]*codeLine{newCodeLine(newLine, true)}, block.lines[insertPos:]...)...)
// 				return
// 			}
// 		}
// 		slog.Error(fmt.Sprintf("Could not insert '%s' into block '%s' after '%s'.", newLine, blockID, insertAfterPrefix))
// 	} else {
// 		slog.Error("No block found with ID: " + blockID)
// 	}
// }

// updating a line that starts with the given prefix, in the block targeted by the given ID;
// if the prefix is not found, a new line is inserted
func (thisCode *codeFile) updateLineIntoBlockWithPrefix(blockID string, newLine string, updateAtPrefix string, insertBeforePrefix string) {
	if block := thisCode.blocksMap[blockID]; block != nil {
		// if the line exists, we update it
		for _, line := range block.lines {
			if line != nil && strings.HasPrefix(line.content, updateAtPrefix) {
				line.rawline = newLine
				return
			}
		}

		// aw, seems like we need to insert it
		if insertBeforePrefix != "" {
			thisCode.insertLineIntoBlockBeforePrefix(blockID, newLine, insertBeforePrefix)
		}
	} else {
		slog.Error("No block found with ID: " + blockID)
	}
}

// tells if the targeted block does have a line with starts with the given prefix, or not
func (thisCode *codeFile) blockHasLineStartingWith(blockID string, prefix string) bool {
	if block := thisCode.blocksMap[blockID]; block != nil {
		for _, line := range block.lines {
			if line != nil && strings.HasPrefix(line.content, prefix) {
				return true
			}
		}
	} else {
		slog.Error("No block found with ID: " + blockID)
	}

	return false
}

// appending a line to the current block, creating the latter if absent
func (thisCode *codeFile) addLineToCurrentBlock(rawline string, isCode bool, content string) {
	if thisCode.current == nil {
		// forcing "isCode = false" here, since a block started this way can't be code
		thisCode.addNewBlock(rawline, false, "", false)
	} else {
		thisCode.current.lines = append(thisCode.current.lines, newCodeLine(rawline, isCode).with(content))
	}
}

// main function reading a file into a structure containing code blocks in a organized fashion
func parseCode(filepath string) (code *codeFile) {
	// what's gonna be returned
	code = &codeFile{
		blocks:    []*codeBlock{},
		blocksMap: make(map[string]*codeBlock),
	}

	// If the file does not exist, the parsing is quite straightforward :)
	if !core.FileExists(filepath) {
		return
	}

	// Opening the file
	file, err := os.Open(filepath)
	if err != nil {
		core.PanicMsgIfErr(err, "Could not open file '%s'", filepath)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Scan the file line by line and append to the slice
	counter := 0     // counting the lines, for debug
	comment := false // capturing /* block \n comments */
	for scanner.Scan() {
		counter += 1
		rawline := scanner.Text()
		content := strings.TrimSpace(rawline)

		// if we're starting or continuing a comment block
		if strings.HasPrefix(content, "/*") || comment {
			code.addLineToCurrentBlock(rawline, false, "")
			comment = !strings.HasSuffix(content, "*/") // if no 1-liner block, the commenting continues
			continue
		}

		// else, we're in "normal code"
		switch {
		case content == "" || strings.HasPrefix(content, "//"):
			code.addLineToCurrentBlock(rawline, false, "")

		case strings.HasPrefix(content, "import"):
			// using the import alias as an ID
			importID := core.IfThenElse(strings.HasPrefix(content, "import * as _"), extractBlockID(content, 12, "from"), "")
			code.addNewBlock(rawline, true, importID, false)

		// case strings.HasPrefix(content, "const") && strings.Contains(content, "fieldAtom("):
		case strings.HasPrefix(rawline, "const"):
			// using what's between "const" and "=" as an ID
			code.addNewBlock(rawline, true, extractBlockID(content, 5, "="), false)

		// case strings.HasPrefix(content, "export const") && strings.Contains(content, "fieldConfigAtom("):
		case strings.HasPrefix(content, "export const"):
			// using what's between "export const" and "=" as an ID
			code.addNewBlock(rawline, true, extractBlockID(content, 12, "="), false)

		case strings.HasPrefix(content, "export function"):
			// using what's between "export const" and "=" as an ID
			code.addNewBlock(rawline, true, extractBlockID(content, 15, "("), false)

		default:
			code.addLineToCurrentBlock(rawline, true, content)
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		core.PanicMsgIfErr(err, "Could not scan file '%s'", filepath)
	}

	return
}

// building an ID from a line, starting from an index, and stopping at a given character
func extractBlockID(content string, start int, end string) string {
	return strings.TrimSpace(content[start:strings.Index(content, end)])
}
