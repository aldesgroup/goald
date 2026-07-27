// ------------------------------------------------------------------------------------------------
// This is about JSON-unmarshalling business objects that have polymorphic relationships (i.e.
// interface-typed fields), which the standard encoding/json package cannot handle on its own,
// since it doesn't know which concrete type to instantiate.
//
// Instead of walking the target struct with the reflect package, we rely on what Goald already
// knows about a business object's shape: its IBusinessObjectModel (built once, at startup) tells
// us which JSON properties are relationships, and the generated, reflection-free
// SetRelationshipValue() / AddRelationshipValue() methods (in each class's "--map.go" file) let us
// assign the resolved target(s) using plain type assertions.
//
// The convention here is that any JSON object meant to fill a polymorphic relationship must carry
// a "class" property, valued with the name of the registered, concrete BO class to use,
// e.g.: "mainContact": {"class": "Employee", "firstName": "John", ...}
// ------------------------------------------------------------------------------------------------
package goald

// TODO TODO TODO: review and master this

import (
	"bytes"
	"encoding/json"

	core "github.com/aldesgroup/corego"
)

// the JSON property expected to carry the concrete class name for a polymorphic relationship
const polymorphicClassField = "class"

// unmarshalling JSON data into a business object, resolving its relationships (if any) using the
// business object's model, and assigning them via the class's generated, reflection-free setters
func unmarshalBObj(data []byte, clsName className, bObj any) error {
	ibObj, isBObj := bObj.(IBusinessObject)
	model := modelForName(clsName)

	// nothing to do here, this isn't a business object we know about: falling back to the standard unmarshalling
	if !isBObj || model == nil || len(model.base().relationships) == 0 {
		return json.Unmarshal(data, bObj)
	}

	// splitting the input JSON into its top-level properties, so we can single out the relationships
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	class := classForName(clsName, true)

	for _, relationship := range model.base().relationships {
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
			if err := class.ClearRelationshipValue(ibObj, relationship.GetName()); err != nil {
				return err
			}

			for _, rawItem := range rawItems {
				target, errTarget := unmarshalRelationshipTarget(relationship, rawItem)
				if errTarget != nil {
					return errTarget
				}
				if err := class.AddRelationshipValue(ibObj, relationship.GetName(), target); err != nil {
					return err
				}
			}
		} else {
			target, errTarget := unmarshalRelationshipTarget(relationship, rawVal)
			if errTarget != nil {
				return errTarget
			}
			if err := class.SetRelationshipValue(ibObj, relationship.GetName(), target); err != nil {
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
// - for a polymorphic relationship, the concrete class is read from the "class" discriminator property
// - for a monomorphic relationship, the concrete class is already known, from the model itself
func unmarshalRelationshipTarget(relationship *Relationship, rawVal json.RawMessage) (IBusinessObject, error) {
	var targetClsName className

	if relationship.IsPolymorphic() {
		var discriminator struct {
			Class string `json:"class"`
		}
		if err := json.Unmarshal(rawVal, &discriminator); err != nil {
			return nil, ErrorC(err, "Could not read the '%s' discriminator property for the '%s' relationship",
				polymorphicClassField, relationship.GetName())
		}
		if discriminator.Class == "" {
			return nil, Error("Missing '%s' property to determine the concrete type to use for the '%s' relationship",
				polymorphicClassField, relationship.GetName())
		}

		targetClsName = className(discriminator.Class)
	} else {
		targetNames := relationship.getTargetClassNames()
		if len(targetNames) == 0 {
			return nil, Error("Relationship '%s' has no target class", relationship.GetName())
		}

		targetClsName = targetNames[0]
	}

	targetClass := classForName(targetClsName, true)

	target := targetClass.NewObject()
	if err := unmarshalBObj(rawVal, targetClsName, target); err != nil {
		return nil, err
	}

	targetBObj, ok := target.(IBusinessObject)
	if !ok {
		return nil, Error("Class '%s' does not describe a business object", targetClsName)
	}

	return targetBObj, nil
}

// returns true if the given raw JSON value is the "null" literal, once trimmed
func isNullJSON(rawVal json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(rawVal), []byte("null"))
}
