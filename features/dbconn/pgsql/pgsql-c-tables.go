package pgsql

import (
	"fmt"
)

// DropTableFkQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) DropTableFkQuery(tableName string, fkName string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s", tableName, fkName)
}

// AddNotNullQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) AddNotNullQuery(tableName string, columnName string, columnType string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s SET NOT NULL", tableName, columnName)
}

// DropNotNullQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) DropNotNullQuery(tableName string, columnName string, columnType string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL", tableName, columnName)
}

// ModifyColumnQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) ModifyColumnQuery(tableName string, columnName string, columnType string) string {
	return fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, columnName, columnType)
}
