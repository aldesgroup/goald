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

	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// constants, variables, useful structs... & main generation function
// ------------------------------------------------------------------------------------------------

var (
	modelsDIRPATH = path.Join("src", "components", "models")
)

const (
	fieldAtomSUFFIX = "FieldAtom"
	formBlockID     = "Form"
)

// TODO maybe do not collocate everything on the server... let's change the receiver here

func (thisServer *server) generateWebAppModels(webdir string, regen bool) {
	// the enum files to generate
	enums := map[string]IEnum{}

	// scanning for BOs involved in endpoints used from the web app
	for _, ep := range restRegistry.endpoints {
		if ep.isFromWebApp() {
			// generating the model for the endpoint resource
			generateWebAppModel(webdir, ep.getResourceClass(), enums, regen)

			// if the endpoint admits a BO as an input (body or URL params), then we also need the model in the webapp
			if inputOrParamsClass := ep.getInputOrParamsClass(); inputOrParamsClass != "" {
				generateWebAppModel(webdir, inputOrParamsClass, enums, regen)
			}
		}
	}

	// enum files generation
	for enumType, enum := range enums {
		// TODO run in go routines
		generateWebAppEnum(webdir, enumType, enum)
	}
}

type codeContext struct {
	enums      map[string]IEnum
	bObjType   *utils.GoaldType
	boInstance *utils.GoaldValue
}

func (ctx *codeContext) getEnumType(field IField) string {
	return ctx.bObjType.FieldByName(field.getName()).Type().Name()
}

func generateWebAppModel(webdir string, clsName className, enums map[string]IEnum, regen bool) {
	// the business object we're dealing with
	boClass := classForName(clsName)
	clUtils := getClassUtils(boClass)
	boFields := utils.GetSortedValues(boClass.base().fields)

	// the file we're dealing with
	filename := utils.PascalToCamel(string(clsName)) + ".ts"
	filepath := path.Join(webdir, modelsDIRPATH, filename)

	// gathering needed info into a context
	codeCtx := &codeContext{
		enums:      enums,
		bObjType:   utils.TypeOf(clUtils.NewObject(), true),
		boInstance: utils.ValueOf(clUtils.NewObject()),
	}

	// gathering the needed enums
	for _, field := range boFields {
		if field.getTypeFamily() == utils.TypeFamilyENUM {
			enums[codeCtx.getEnumType(field)] = codeCtx.boInstance.GetFieldValue(field.getName()).(IEnum)
		}
	}

	// TODO remove
	if !regen && utils.FileExists(filepath) && utils.EnsureModTime(filepath).After(clUtils.getLastBOMod()) {
		return // the file already exists and is older than our changes in the BO class file
	}

	// getting the file content - which might be empty if the file does not exist yet
	code := parseCode(filepath).addImportsIfNeeded()

	// browsing the entity's properties to fill the get / set cases in the 2 switch
	for _, field := range boFields {
		code.addFieldIfNeeded(codeCtx, field)
	}

	// "unpacking" the code blocks to code lins
	codeLines := []string{}
	for _, block := range code.blocks {
		// the block should at least have a non-nil code line
		if block.lines[0] != nil {
			// codeLines = append(codeLines, fmt.Sprintf("// --- block [%s] ---------------------------", block.id))
			for _, line := range block.lines {
				if line != nil {
					// codeLines = append(codeLines, line.content+fmt.Sprintf(" // %d", line.num))
					codeLines = append(codeLines, line.rawline)
				}
			}
		}
	}

	// writing out the code lines
	utils.WriteToFile(strings.Join(codeLines, newline), filepath)
	slog.Info(fmt.Sprintf("(Re-)generated file %s", filepath))
}

// ------------------------------------------------------------------------------------------------
// code editing
// ------------------------------------------------------------------------------------------------

// initiating the code
func (thisCode *codeFile) addImportsIfNeeded() *codeFile {
	if len(thisCode.blocks) == 0 {
		thisCode.addNewBlock("import { fieldAtom, formAtom, useFieldActions, useFieldValue, useInputField } from \"form-atoms\";", true, "", false)
		thisCode.addNewBlock("import { atom, useSetAtom } from \"jotai\";", true, "", false)
		thisCode.addNewBlock("import { useEffect } from \"react\";", true, "", false)
		thisCode.addNewBlock("import { fieldConfigAtom } from \"~/vendor/goaldr\";", true, "", false)
		thisCode.addNewBlock("export const "+formBlockID+" = formAtom({", true, formBlockID, true).appendLine("});", true)

	}

	return thisCode
}

// handling a field, adding it if not in the code already, flagging an enum for generation if it's an enum field
func (thisCode *codeFile) addFieldIfNeeded(codeCtx *codeContext, field IField) {
	// adding to the context, and the class file content
	if typeFamily := field.getTypeFamily(); typeFamily != utils.TypeFamilyUNKNOWN && typeFamily != utils.TypeFamilyRELATIONSHIP {
		// not handling multiple properties for now - nor the ID field
		if !field.isMultiple() && field.getName() != "ID" {
			var (
				enumType, enumVar, initVal string
			)

			// dealing with some field specificities
			switch typeFamily {
			// --- enums -------------------------------------------------------------------
			case utils.TypeFamilyENUM:
				// flagging this enum type for code generation
				enumType = codeCtx.getEnumType(field)

				// importing the enum if needed
				enumVar = "_" + utils.PascalToCamel(enumType)
				if thisCode.blocksMap[enumVar] == nil {
					// panic(fmt.Sprintf("Missing import for '%s'; should add it after block %d", enumType, thisCode.findLastImportPosition()))
					thisCode.addEnumImport(enumVar)
				}

				// proposing an init value
				initVal = fmt.Sprintf("%s.%s", enumVar, makeEnumName(utils.GetFirstMapValue(codeCtx.enums[enumType].Values())))

			// --- numbers -----------------------------------------------------------------
			case utils.TypeFamilyINT, utils.TypeFamilyBIGINT, utils.TypeFamilyREAL, utils.TypeFamilyDOUBLE:
				// proposing an init value
				initVal = "0"

			// --- booleans ----------------------------------------------------------------
			case utils.TypeFamilyBOOL:
				// proposing an init value
				initVal = "false"
			}

			// adding the field atom declaration if needed
			fieldAtomName := field.getName() + fieldAtomSUFFIX
			if thisCode.blocksMap[fieldAtomName] == nil {
				fieldAtomDecl := fmt.Sprintf("const %s = fieldAtom({ value: %s? });", fieldAtomName, initVal)
				thisCode.addNewBlockBeforeLast(fieldAtomDecl, true, fieldAtomName, true)
			}

			// adding the field atom to the form if needed
			if !thisCode.blockHasLineStartingWith(formBlockID, field.getName()+":") {
				thisCode.insertLineIntoBlockBeforePrefix(formBlockID, fmt.Sprintf("    %s: %s,", field.getName(), fieldAtomName), "}")
			}

			// adding the field config atom if needed
			missingConfigAtom := thisCode.blocksMap[field.getName()] == nil
			if missingConfigAtom {
				fieldConfigAtomDecl := fmt.Sprintf("export const %s = fieldConfigAtom({", field.getName())
				newBlock := thisCode.addNewBlockBeforeLast(fieldConfigAtomDecl, true, field.getName(), false)
				newBlock.appendLine(fmt.Sprintf("    fieldAtom: %s,", fieldAtomName), true)
				newBlock.appendLine("});", true)

			}

			// linking the enum's options to the field, if needed
			if (typeFamily == utils.TypeFamilyENUM) && (missingConfigAtom || !thisCode.blockHasLineStartingWith(field.getName(), "options:")) {
				thisCode.insertLineIntoBlockBeforePrefix(field.getName(), fmt.Sprintf("    options: %s.Options,", enumVar), "}")
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
	for pos = 3; pos < len(thisCode.blocks); pos++ {
		if !strings.HasPrefix(thisCode.blocks[pos].lines[0].rawline, "import ") {
			break
		}
	}

	return
}

// ------------------------------------------------------------------------------------------------
// enum files generation
// ------------------------------------------------------------------------------------------------

func generateWebAppEnum(webdir string, enumType string, enum IEnum) {
	filename := "_" + utils.PascalToCamel(enumType) + ".ts"
	filepath := path.Join(webdir, modelsDIRPATH, filename)

	content := ""
	allTypes := []string{}
	for _, enumVal := range utils.GetSortedKeys(enum.Values()) {
		enumLabel := enum.Values()[enumVal]
		enumName := makeEnumName(enumLabel)
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

// ------------------------------------------------------------------------------------------------
// code parsing & generic code editing
// ------------------------------------------------------------------------------------------------

type codeLine struct {
	rawline string
	isCode  bool
	content string
	// num     int
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
		// num:     num,
	}
}

type codeBlock struct {
	// we should be able to insert new lines, for new min / max constraints for instance, or new field atoms in the form
	lines []*codeLine
	id    string
}

// how to add a line to a code block
func (block *codeBlock) appendLine(rawline string, isCode bool) {
	block.lines = append(block.lines, newCodeLine(rawline, isCode))
}

type codeFile struct {
	blocks    []*codeBlock
	blocksMap map[string]*codeBlock // some code blocks mapped with their ID
	current   *codeBlock            // used during the parsing
	// counter   int                   // line counter
}

// adding a new block of code after the last one already present
func (thisCode *codeFile) addNewBlock(rawline string, isCode bool, id string, blankBefore bool) *codeBlock {
	return thisCode.addNewBlockAtPosition(rawline, isCode, id, blankBefore, -1)
}

// adding a new block of code before the last one already present
func (thisCode *codeFile) addNewBlockBeforeLast(rawline string, isCode bool, id string, blankBefore bool) *codeBlock {
	return thisCode.addNewBlockAtPosition(rawline, isCode, id, blankBefore, len(thisCode.blocks)-1)
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
		// thisCode.counter++
		// thisCode.current.lines = append(thisCode.current.lines, newCodeLine("", thisCode.counter, isCode)) TODO remove TODO change
		thisCode.current.appendLine("", isCode)
	}

	// adding the code line
	// thisCode.counter++
	// thisCode.current.lines = append(thisCode.current.lines, newCodeLine(line, thisCode.counter, isCode)) TODO remove TODO change
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

// appending a line before the last one to the block target by the given ID
func (thisCode *codeFile) insertLineIntoBlockBeforePrefix(blockID string, newLine string, prefix string) {
	if block := thisCode.blocksMap[blockID]; block != nil {
		// finding where to cut the block lines
		var pos int
		for pos = len(block.lines) - 1; pos >= 0; pos-- {
			line := block.lines[pos]
			if line != nil && strings.HasPrefix(line.rawline, prefix) {
				break
			}
		}

		// inserting the new line at the right position
		block.lines = append(block.lines[:pos], append([]*codeLine{newCodeLine(newLine, true)}, block.lines[pos:]...)...)
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
		// thisCode.counter++
		// thisCode.current.lines = append(thisCode.current.lines, newCodeLine(line, thisCode.counter, isCode))
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
	if !utils.FileExists(filepath) {
		return
	}

	// Opening the file
	file, err := os.Open(filepath)
	if err != nil {
		utils.PanicErrf(err, "Could not open file '%s'", filepath)
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
			importID := utils.IfThenElse(strings.HasPrefix(content, "import * as _"), extractBlockID(content, 12, "from"), "")
			code.addNewBlock(rawline, true, importID, false)

		// case strings.HasPrefix(content, "const") && strings.Contains(content, "fieldAtom("):
		case strings.HasPrefix(rawline, "const"):
			// using what's between "const" and "=" as an ID
			code.addNewBlock(rawline, true, extractBlockID(content, 5, "="), false)

		// case strings.HasPrefix(content, "export const") && strings.Contains(content, "fieldConfigAtom("):
		case strings.HasPrefix(content, "export const"):
			// using what's between "export const" and "=" as an ID
			code.addNewBlock(rawline, true, extractBlockID(content, 12, "="), false)

		default:
			code.addLineToCurrentBlock(rawline, true, content)
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		utils.PanicErrf(err, "Could not scan file '%s'", filepath)
	}

	return
}

// building an ID from a line, starting from an index, and stopping at a given character
func extractBlockID(content string, start int, end string) string {
	return strings.TrimSpace(content[start:strings.Index(content, end)])
}
