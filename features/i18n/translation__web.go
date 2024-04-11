package i18n

import (
	"fmt"

	g "github.com/aldesgroup/goald"
	// "github.com/aldesgroup/goald/_generated/class"
	"github.com/aldesgroup/goald/features/hstatus"
)

func init() {
	g.GetManyWithParams[*Translation, *TranslationUrlParams](listTranslations, "").
		// TargetWith(class.Translation().LangStr()).
		Label("Returns the translations for the given route")
}

func listTranslations(webCtx g.WebContext, params *TranslationUrlParams) ([]*Translation, hstatus.Code, string) {
	// getting the targeted language
	langStr := webCtx.GetTargetRefOrID()

	// getting the translations for the right language
	translations, errGet := getTranslations(webCtx.GetBloContext(), langStr, params.Route, params.Part, params.Key)
	if errGet != nil {
		return nil, hstatus.InternalServerError, fmt.Sprintf("Error while searching for translations: %s", errGet)
	}
	if len(translations) == 0 {
		return nil, hstatus.NotFound, fmt.Sprintf("No translation found for lang = %s / route = %s / part = %s / key = %s",
			langStr, params.Route, params.Part, params.Key)
	}

	return translations, hstatus.OK, ""
}
