// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

// checking a business object's general validity
func (thisClass *UserGroupClass) IsModelValid(bObj goald.IBusinessObject) error {
	bo := bObj.(*accessmgt.UserGroup)

	if err := goald.CheckStringSize(bo.Description, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Description'")
	}
	if bo.Name == "" {
		return goald.Error("'Name' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Name, 24, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Name'")
	}

	return nil
}
