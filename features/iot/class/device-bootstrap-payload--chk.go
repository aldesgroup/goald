// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// checking a business object's general validity
func (thisClass *DeviceBootstrapPayloadClass) IsModelValid(bObj goald.IBusinessObject) error {
	bo := bObj.(*iot.DeviceBootstrapPayload)

	if bo.Model == "" {
		return goald.Error("'Model' is mandatory and must have a non-zero value")
	}
	if bo.Nonce == "" {
		return goald.Error("'Nonce' is mandatory and must have a non-zero value")
	}
	if bo.Serial == "" {
		return goald.Error("'Serial' is mandatory and must have a non-zero value")
	}
	if bo.Timestamp == 0 {
		return goald.Error("'Timestamp' is mandatory and must have a non-zero value")
	}

	return nil
}
