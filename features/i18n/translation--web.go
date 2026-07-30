package i18n

import (
	"fmt"

	g "github.com/aldesgroup/goald"
	// "github.com/aldesgroup/goald/_include/i18n/model"
	"github.com/aldesgroup/goald/_include/i18n/model"
	"github.com/aldesgroup/goald/features/hstatus"
)

var GenericEPs = g.NewEndpointGroup("General", "Generic endpoints for common behaviors")

func init() {
	g.GetManyWithParams[*Translation, *TranslationUrlParams](listTranslations, "").
		InGroup(GenericEPs).
		TargetWith(model.Translation().Lang()).
		Label("List translations").
		Description("Returns the translations for the given route")
}

func listTranslations(webCtx g.WebContext, params *TranslationUrlParams) ([]*Translation, hstatus.Code, string) {
	// getting the targeted language
	langStr := webCtx.GetResourceRefOrID()

	// not translating english - FOR NOW // TODO refine
	if langStr == "en" {
		return nil, hstatus.OK, ""
	}

	// getting the translations for the right language
	foundTranslations, errGet := getTranslations(webCtx.GetBloContext(), langStr, params.Namespace, params.Key)
	if errGet != nil {
		return nil, hstatus.InternalServerError, fmt.Sprintf("Error while searching for translations: %s", errGet)
	}
	if len(foundTranslations) == 0 {
		return nil, hstatus.NotFound, fmt.Sprintf("No translation found for lang = %s / namespace = %s / key = %s",
			langStr, params.Namespace, params.Key)
	}

	return foundTranslations, hstatus.OK, ""
}
