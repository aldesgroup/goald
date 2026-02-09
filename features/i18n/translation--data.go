// ------------------------------------------------------------------------------------------------
// Here, mainly about loading the translations at startup
// ------------------------------------------------------------------------------------------------
package i18n

import (
	"fmt"
	"log/slog"
	"path"
	"strings"
	"sync"

	core "github.com/aldesgroup/corego"
	g "github.com/aldesgroup/goald"
)

// ------------------------------------------------------------------------------------------------
// Making this data loader function available
// ------------------------------------------------------------------------------------------------

func init() {
	g.RegisterDataLoader(loadTranslations, false)
}

// ------------------------------------------------------------------------------------------------
// Useful variables / constants / structs
// ------------------------------------------------------------------------------------------------

const fileLoadingWorkersNb = 5 // we've computed with a benchmark that the optimum lies around this value

type fileLoadingWorkerCtx struct {
	lang         Language
	file         string
	ns           string
	translations []*Translation
}

func (ctx *fileLoadingWorkerCtx) namespace() string {
	if ctx.ns == "" {
		ctx.ns = strings.TrimSuffix(ctx.file, ".json")
	}

	return ctx.ns
}

// ------------------------------------------------------------------------------------------------
// Main data loading function & utils
// ------------------------------------------------------------------------------------------------

func loadTranslations(ctx g.BloContext, params map[string]string) error {
	// checking the parameters
	if params == nil {
		return g.Error("No 'loadTranslations' data loader item in the config!")
	}
	folderPath := params["folder"]
	if folderPath == "" {
		return g.Error("Empty value for 'loadTranslations.folder' in the config!")
	}
	if !core.DirExists(folderPath) {
		return g.Error("'%s' is not a valid path!", folderPath)
	}

	// finding out all the languages, and all the translation files (1 file per namespace)
	languages := []Language{}
	filesNb := 0
	files := map[Language][]string{}
	for _, languageEntry := range core.EnsureReadDir(folderPath) {
		// we should have 1 folder per language
		if languageEntry.IsDir() {
			// there we have our language
			language := LanguageFrom(languageEntry.Name())
			languages = append(languages, language)

			// let's gather the translation files
			for _, fileEntry := range core.EnsureReadDir(path.Join(folderPath, languageEntry.Name())) {
				filesForLanguage := []string{}
				if strings.HasSuffix(fileEntry.Name(), ".json") {
					filesNb++
					filesForLanguage = append(filesForLanguage, fileEntry.Name())
				}
				files[language] = append(files[language], filesForLanguage...)
			}
		}
	}

	slog.Debug(fmt.Sprintf("Found %d languages and %d translation files in total", len(languages), filesNb))

	// prepping for loading of all the translation files done in parallel
	filesToLoad := make(chan *fileLoadingWorkerCtx, (filesNb))
	loadedFiles := make(chan *fileLoadingWorkerCtx, (filesNb))
	waitGroup := new(sync.WaitGroup)

	// starting workers to deal with the files to load
	for workerID := 0; workerID < fileLoadingWorkersNb; workerID++ {
		waitGroup.Go(func() {
			for fileToLoad := range filesToLoad {
				doLoadFile(workerID, folderPath, fileToLoad, loadedFiles)
			}
		})
	}

	// adding files to load
	for language, filesForLanguage := range files {
		for _, file := range filesForLanguage {
			filesToLoad <- &fileLoadingWorkerCtx{language, file, "", nil}
		}
	}

	// we're done adding stuff to do
	close(filesToLoad)

	// waiting for all the workers to be done
	waitGroup.Wait()

	// we won't be getting new
	close(loadedFiles)

	// (re-)initialising the global object containing all the translations
	allTranslations = make(map[Language]map[string][]*Translation)

	// transfering the translations from the loaded files to our "big" translations map
	for loadedFile := range loadedFiles {
		// initialising for the language, if needed
		if allTranslations[loadedFile.lang] == nil {
			allTranslations[loadedFile.lang] = make(map[string][]*Translation)
		}

		// now adding the translations for the current language & namespace
		allTranslations[loadedFile.lang][loadedFile.namespace()] = loadedFile.translations
	}

	// TODO ?
	// debounce the restart of containers when aldev refresh runs OR use another place for the compilation / run
	// hide the logs for the server start, only show in verbose mode
	// aldev should also have a silent mode by default

	return nil
}

func doLoadFile(workerID int, folderPath string, fileToLoad *fileLoadingWorkerCtx, loadedFiles chan *fileLoadingWorkerCtx) {
	// the path of the file to load
	filePath := path.Join(folderPath, fileToLoad.lang.String(), fileToLoad.file)

	// a bit of logging
	// slog.Debug(fmt.Sprintf("Worker %02d: loading '%s'", workerID, filePath))

	// JSON => key-value map of the translations
	rawTranslations := *core.ReadFileFromJSON(filePath, &map[string]string{}, true)

	// building Translation objects from this map
	for _, key := range core.GetSortedKeys(rawTranslations) {
		fileToLoad.translations = append(fileToLoad.translations, &Translation{
			Lang:      fileToLoad.lang.String(),
			Namespace: fileToLoad.namespace(),
			Key:       key,
			Value:     rawTranslations[key],
		})
	}

	// that's one more file completely loaded!
	loadedFiles <- fileToLoad
}
