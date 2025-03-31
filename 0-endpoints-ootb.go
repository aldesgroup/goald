// ------------------------------------------------------------------------------------------------
// Out-of-the-box endpoints to quickly have SCRUD working with Business Objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"github.com/aldesgroup/goald/features/hstatus"
)

// Simply listing the resources of a targeted type
func HandleGetAll[ResourceType IBusinessObject](webCtx WebContext) ([]ResourceType, hstatus.Code, string) {
	list, errList := LoadBOs[ResourceType](webCtx.GetBloContext(), webCtx.GetResource(), webCtx.GetResourceLoadingType())
	if errList != nil {
		msg := ErrorC(errList, "Could not get a list of '%s' instances", webCtx.GetResource().base().name).Error()
		return nil, hstatus.InternalServerError, msg
	}

	return list, hstatus.OK, ""
}
