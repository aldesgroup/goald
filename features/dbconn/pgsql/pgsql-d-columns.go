package pgsql

import (
	"fmt"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
)

// SQLColumnDeclaration implements [goald.iDBAdapter].
// It returns the SQL column declaration for the given BO property, including its type and a potential NOT NULL constraint.
// For polymorphic relationships, it returns a column dec laration for the foreign key column, plus a column for the target object type.
func (thisAdapter *dbAdapterPGSQL) SQLColumnDeclaration(property goald.IBusinessObjectProperty) (string, string) {
	if property.GetName() == goald.BoFieldID {
		return "BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY", ""
	}

	notNull := ""
	if property.IsRequiredInDb() {
		notNull = " NOT NULL"
	}

	// cast on the actual type of the property, to access its specific methods
	switch property := property.(type) {

	case *goald.Relationship:
		if property.IsPolymorphic() {
			// this is the special case, where we need to add an additional column for the target object type
			return "BIGINT" + notNull, "VARCHAR(48)"
		}
		return "BIGINT" + notNull, ""

	case *goald.IntField:
		return "INTEGER" + notNull, ""

	case *goald.BigIntField:
		return "BIGINT" + notNull, ""

	case *goald.StringField:
		return fmt.Sprintf("VARCHAR(%d)", property.GetSize()) + notNull, ""

	case *goald.RealField:
		return fmt.Sprintf("NUMERIC(%d, %d)", property.GetTotalDigits(), property.GetDecimals()) + notNull, ""

	case *goald.DoubleField:
		return fmt.Sprintf("NUMERIC(%d, %d)", property.GetTotalDigits(), property.GetDecimals()) + notNull, ""

	case *goald.BoolField:
		return "BOOLEAN NOT NULL DEFAULT FALSE", ""

	case *goald.DateField:
		return "TIMESTAMPTZ" + notNull, ""

	case *goald.EnumField:
		if !property.IsMultiple() {
			return "INTEGER" + notNull, ""
		}

	default:
		core.PanicMsg("Cannot build the SQL declaration for property '%s' (type: %T, multiple: %v)", property.GetName(), property, property.IsMultiple())
	}

	// this should never happen
	return "- NOT HANDLED YET -", ""
}
