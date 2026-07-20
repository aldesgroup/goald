// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

// checking a business object's general validity
func (thisClass *TranslationClass) IsModelValid(bObj goald.IBusinessObject) error {
	bo := bObj.(*i18n.Translation)

	if bo.Lang == "" {
		return goald.Error("'Lang' is mandatory and must have a non-zero value")
	}

	return nil
}
