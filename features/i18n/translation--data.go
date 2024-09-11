// ------------------------------------------------------------------------------------------------
// Here, mainly about loading the translations at startup
// ------------------------------------------------------------------------------------------------
package i18n

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	g "github.com/aldesgroup/goald"
)

func init() {
	g.RegisterDataLoader(loadTranslations, false)
}

// the structure of the translation file as downloaded by Aldev
type translationRow struct {
	Namespace string   `json:"n"`
	Key       string   `json:"k"`
	Values    []string `json:"v"`
}

func loadTranslations(ctx g.BloContext, params map[string]string) error {
	slog.Info("Loading the translations...")

	if params == nil {
		return g.Error("No 'loadTranslations' data loader item in the config!")
	}
	if params["file"] == "" {
		return g.Error("Empty value for 'loadTranslations.file' in the config!")
	}

	// reading the translation file
	dataBytes, errRead := os.ReadFile(params["file"])
	if errRead != nil {
		return g.ErrorC(errRead, "Could not read file '%s'", params["file"])
	}

	// unmarshaling the translation rows found in the file
	translationRows := []*translationRow{}
	if errUnmarsh := json.Unmarshal(dataBytes, &translationRows); errUnmarsh != nil {
		return g.ErrorC(errUnmarsh, "Could not unmarshal the translation data")
	}

	// initialising the global object containing all the translations
	allTranslations = make(map[Language]map[string][]*Translation)

	// initialising each language
	for _, translation := range translationRows[0].Values {
		lang, _ := getLangAndValue(translation)
		if lang == LanguageUNDEFINED {
			return g.Error("Undefined language found here: %s", translation[:2])
		}
		allTranslations[lang] = make(map[string][]*Translation)
	}

	// "registering" the translations now
	for _, row := range translationRows {
		for _, translation := range row.Values {
			// getting the current language + translation value in this language
			lang, value := getLangAndValue(translation)

			// adding the value to the right route
			allTranslations[lang][row.Namespace] = append(allTranslations[lang][row.Namespace], &Translation{
				Namespace: row.Namespace,
				Key:       row.Key,
				Value:     value,
			})
		}
	}

	// TODO
	// debounce the restart of containers when aldev complete runs OR use another place for the compilation / run
	// hide the logs for the server start, only show in verbose mode
	// aldev should also have a silent mode by default

	return nil
}

func getLangAndValue(translation string) (Language, string) {
	langAndValue := strings.SplitN(translation, ": ", 2)
	return LanguageFrom(langAndValue[0]), langAndValue[1]
}
