package i18n

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// constants & variables
// ------------------------------------------------------------------------------------------------

const AnyVALUE = "any"

// the translations are mapped by language, then route
var allTranslations map[Language]map[string][]*Translation

// ------------------------------------------------------------------------------------------------
// the different languages Goald applications can manage
// ------------------------------------------------------------------------------------------------

func getTranslations(_ g.BloContext, langStr, route, partArg, key string) ([]*Translation, error) {
	if allTranslations == nil {
		return nil, g.Error("The translations have not been loaded!")
	}

	// the language has to be known
	lang := LanguageFrom(langStr)
	if lang == LanguageUNDEFINED {
		return nil, g.Error("Unhandled language: '%s'", langStr)
	}

	// the route cannot be empty - at least for now
	if route == "" {
		return nil, g.Error("The 'Route' has to be provided!")
	}

	part := utils.IfThenElse(partArg != "", partArg, AnyVALUE)

	translationsForLang := allTranslations[lang]
	if len(translationsForLang) == 0 {
		return nil, nil
	}

	translationsForRoute := translationsForLang[route]
	if len(translationsForRoute) == 0 {
		return nil, nil
	}

	// we don't care about a specific part or key, so all the chosen route's translations it is!
	if part == AnyVALUE && key == "" {
		return translationsForRoute, nil
	}

	// else, we have to filter on the part & maybe the key
	specificTranslations := []*Translation{}
	for _, translation := range translationsForRoute {
		if translation.Part != part {
			continue
		}
		if key != "" && translation.Key != key {
			continue
		}

		// no filtering here => adding
		specificTranslations = append(specificTranslations, translation)
	}

	return specificTranslations, nil
}
