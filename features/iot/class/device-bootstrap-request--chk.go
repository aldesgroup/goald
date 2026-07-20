// Generated file, do not edit!
package class

import (
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// checking a business object's general validity
func (thisClass *DeviceBootstrapRequestClass) IsModelValid(bObj goald.IBusinessObject) error {
	bo := bObj.(*iot.DeviceBootstrapRequest)

	if bo.FactoryCertPEM == "" {
		return goald.Error("'FactoryCertPEM' is mandatory and must have a non-zero value")
	}
	if bo.PayloadB64 == "" {
		return goald.Error("'PayloadB64' is mandatory and must have a non-zero value")
	}
	if bo.SignatureB64 == "" {
		return goald.Error("'SignatureB64' is mandatory and must have a non-zero value")
	}

	return nil
}
