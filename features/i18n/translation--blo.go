package i18n

import (
	g "github.com/aldesgroup/goald"
)

// ------------------------------------------------------------------------------------------------
// constants & variables
// ------------------------------------------------------------------------------------------------

const AnyVALUE = "any"

// the translations are mapped by language, then namespace
var allTranslations map[Language]map[string][]*Translation

// ------------------------------------------------------------------------------------------------
// the different languages Goald applications can manage
// ------------------------------------------------------------------------------------------------

func getTranslations(_ g.BloContext, langStr, namespace, key string) ([]*Translation, error) {
	if allTranslations == nil {
		return nil, g.Error("The translations have not been loaded!")
	}

	// the language has to be known
	lang := LanguageFrom(langStr)
	if lang == LanguageUNDEFINED {
		return nil, g.Error("Unhandled language: '%s'", langStr)
	}

	// the namespace cannot be empty - at least for now
	if namespace == "" {
		return nil, g.Error("The 'Namespace' has to be provided!")
	}

	translationsForLang := allTranslations[lang]
	if len(translationsForLang) == 0 {
		return nil, nil
	}

	translationsForNamespace := translationsForLang[namespace]
	if len(translationsForNamespace) == 0 {
		return nil, nil
	}

	// we don't care about a specific key, so all the chosen language & namespace's translations it is!
	if key == "" {
		return translationsForNamespace, nil
	}

	// else, we have to filter on the part & maybe the key
	specificTranslations := []*Translation{}
	for _, translation := range translationsForNamespace {
		if key != "" && translation.Key != key {
			continue
		}

		// no filtering here => adding
		specificTranslations = append(specificTranslations, translation)
	}

	return specificTranslations, nil
}
