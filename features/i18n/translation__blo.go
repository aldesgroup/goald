package i18n

import (
	g "github.com/aldesgroup/goald"
)

// ------------------------------------------------------------------------------------------------
// constants & variables
// ------------------------------------------------------------------------------------------------

const AnyVALUE = "any"

// the translations are mapped by language, then route
var translations map[Language]map[string][]*Translation

// ------------------------------------------------------------------------------------------------
// the different languages Goald applications can manage
// ------------------------------------------------------------------------------------------------

func getTranslations(_ g.BloContext, langStr, route, part, key string) ([]*Translation, error) {
	if translations == nil {
		return nil, g.Error("The translations have not been loaded!")
	}

	lang := LanguageFrom(langStr)
	if lang == LanguageUNDEFINED {
		return nil, g.Error("Unhandled language: '%s'", langStr)
	}

	if translationsForLang := translations[lang]; len(translationsForLang) > 0 {
		if translationsForRoute := translationsForLang[route]; len(translationsForRoute) > 0 {
			if part == "" && key == "" {
				return translationsForRoute, nil
			}

		}
	}

	return nil, nil
}
