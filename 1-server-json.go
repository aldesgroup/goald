// ------------------------------------------------------------------------------------------------
// This is about JSON-unmarshalling business objects that have polymorphic relationships (i.e.
// interface-typed fields), which the standard encoding/json package cannot handle on its own,
// since it doesn't know which concrete type to instantiate.
//
// Instead of walking the target struct with the reflect package, we rely on what Goald already
// knows about a business object's shape: its IBusinessObjectModel (built once, at startup) tells
// us which JSON properties are relationships, and the generated, reflection-free
// SetRelationshipValue() / AddRelationshipValue() methods (in each model's "--xtd.go" file) let us
// assign the resolved target(s) using plain type assertions.
//
// The convention here is that any JSON object meant to fill a polymorphic relationship must carry
// a "mdl" property, valued with the name of the registered, concrete BO model to use,
// e.g.: "mainContact": {"mdl": "Employee", "firstName": "John", ...}
// ------------------------------------------------------------------------------------------------
package goald

// TODO TODO TODO: review and master this

import (
	"bytes"
	"encoding/json"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// the JSON property expected to carry the concrete model name for a polymorphic relationship
const polymorphicModelField = "mdl"

// unmarshalling JSON data into a business object, resolving its relationships (if any) using the
// business object's model, and assigning them via the model's generated, reflection-free setters
func unmarshalBObj(data []byte, bObj any) error {
	ibObj, isBObj := bObj.(IBusinessObject)
	model := ibObj.getModel(ibObj)

	// nothing to do here, this isn't a business object we know about: falling back to the standard unmarshalling
	if !isBObj || model == nil || len(model.getRelationships()) == 0 {
		return json.Unmarshal(data, bObj)
	}

	// splitting the input JSON into its top-level properties, so we can single out the relationships
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, relationship := range model.getRelationships() {
		jsonName := core.PascalToCamel(relationship.GetName())

		rawVal, present := raw[jsonName]
		if !present || isNullJSON(rawVal) {
			continue
		}

		// removing this property, so the final, standard unmarshalling below doesn't choke on it
		delete(raw, jsonName)

		if relationship.IsMultiple() {
			var rawItems []json.RawMessage
			if err := json.Unmarshal(rawVal, &rawItems); err != nil {
				return ErrorC(err, "Could not read the '%s' relationship", jsonName)
			}

			// resetting the relationship first, so a provided (possibly empty) array always
			// replaces whatever the target BO already had, rather than appending to it
			if err := ibObj.ClearRelationshipValue(relationship.GetName()); err != nil {
				return err
			}

			for _, rawItem := range rawItems {
				target, errTarget := unmarshalRelationshipTarget(relationship, rawItem)
				if errTarget != nil {
					return errTarget
				}
				if err := ibObj.AddRelationshipValue(relationship.GetName(), target); err != nil {
					return err
				}
			}
		} else {
			target, errTarget := unmarshalRelationshipTarget(relationship, rawVal)
			if errTarget != nil {
				return errTarget
			}
			if err := ibObj.SetRelationshipValue(relationship.GetName(), target); err != nil {
				return err
			}
		}
	}

	// unmarshalling the remaining, plain properties the standard way
	remaining, errMarshal := json.Marshal(raw)
	if errMarshal != nil {
		return errMarshal
	}

	return json.Unmarshal(remaining, bObj)
}

// resolving & instantiating the concrete business object targeted by a relationship's raw JSON value:
// - for a polymorphic relationship, the concrete model is read from the "model" discriminator property
// - for a monomorphic relationship, the concrete model is already known, from the model itself
func unmarshalRelationshipTarget(relationship *Relationship, rawVal json.RawMessage) (IBusinessObject, error) {
	var targetModelName utils.ModelName

	if relationship.IsPolymorphic() {
		var discriminator struct {
			Mdl string `json:"mdl"`
		}
		if err := json.Unmarshal(rawVal, &discriminator); err != nil {
			return nil, ErrorC(err, "Could not read the '%s' discriminator property for the '%s' relationship",
				polymorphicModelField, relationship.GetName())
		}
		if discriminator.Mdl == "" {
			return nil, Error("Missing '%s' property to determine the concrete type to use for the '%s' relationship",
				polymorphicModelField, relationship.GetName())
		}

		targetModelName = utils.ModelName(discriminator.Mdl)
	} else {
		targetNames := relationship.getTargetModelNames()
		if len(targetNames) == 0 {
			return nil, Error("Relationship '%s' has no target model", relationship.GetName())
		}

		targetModelName = targetNames[0]
	}

	targetModel := modelFor(targetModelName, true)

	target := targetModel.NewObject()
	if err := unmarshalBObj(rawVal, target); err != nil {
		return nil, err
	}

	targetBObj, ok := target.(IBusinessObject)
	if !ok {
		return nil, Error("Model '%s' does not describe a business object", targetModelName)
	}

	return targetBObj, nil
}

// returns true if the given raw JSON value is the "null" literal, once trimmed
func isNullJSON(rawVal json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(rawVal), []byte("null"))
}
