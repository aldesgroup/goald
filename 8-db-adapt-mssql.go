package goald

import (
	"fmt"
	"log/slog"

	"github.com/aldesgroup/goald/features/utils"
)

// specific queries for SQL Server databases
type dbAdapterMSSQL struct{}

// checking the compliance with the interface
var _ iDBAdapter = (*dbAdapterMSSQL)(nil)

func (thisAdapter *dbAdapterMSSQL) getConnectionString(conf *dbConfig) string {
	return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s;database=%s;dial timeout=5;connection timeout=5", conf.DbHost, conf.DbPort, conf.User, conf.Password, conf.DbName)
	// return fmt.Sprintf("server=%s;port=%d;user id=%s;password=%s", conf.DbHost, conf.DbPort, conf.User, conf.Password)
}

func (thisAdapter *dbAdapterMSSQL) getTablesQuery(dbName string) string {
	return fmt.Sprintf("SELECT name from %s.sys.tables", dbName)
}

// getSQLColumnDeclaration returns the type of the column to create for the given BO property
func (thisAdapter *dbAdapterMSSQL) getSQLColumnDeclaration(property iBusinessObjectProperty) string {
	notNull := utils.IfThenElse(property.isMandatory(), " NOT NULL", "")

	switch property := property.(type) {
	case *Relationship:
		return property.getColumnName() + " BIGINT" + notNull
	case *BoolField:
		return property.getColumnName() + " BIT"
	case *StringField:
		return property.getColumnName() + fmt.Sprintf(" VARCHAR(%d)", property.size) + notNull
	case *IntField:
		return property.getColumnName() + " INT" + notNull
	case *BigIntField:
		return property.getColumnName() + " BIGINT" + notNull
	case *RealField:
		return property.getColumnName() + " REAL" + notNull
	case *DoubleField:
		return property.getColumnName() + " FLOAT" + notNull
	case *DateField:
		return property.getColumnName() + " DATETIME2(6)" + notNull
	case *EnumField:
		if property.isMultiple() {
			panic("not handling listenums yet!!!")
		}
		return property.getColumnName() + " INT" + notNull
	}

	slog.Error(fmt.Sprintf("Not handling this property in DB: %s", property.getName()))

	return ""
}
