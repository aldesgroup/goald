// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/accessmgt"
)

// checking a business object's general validity
func (thisClass *UserClass) IsModelValid(bObj goald.IBusinessObject) error {
	bo := bObj.(*accessmgt.User)

	if bo.Email == "" {
		return goald.Error("'Email' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Email, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Email'")
	}
	if bo.FirstName == "" {
		return goald.Error("'FirstName' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.FirstName, 24, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'FirstName'")
	}
	if bo.LastName == "" {
		return goald.Error("'LastName' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.LastName, 24, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'LastName'")
	}
	if bo.Password == "" {
		return goald.Error("'Password' is mandatory and must have a non-zero value")
	}
	if err := goald.CheckStringSize(bo.Password, 64, 0); err != nil {
		return goald.ErrorC(err, "Invalid value for 'Password'")
	}

	return nil
}
