package goald

import "fmt"

// specific queries for SQL Server databases
type dbAdapterMSSQL struct{}

func (thisAdapter *dbAdapterMSSQL) getTablesQuery(dbName string) string {
	return fmt.Sprintf("SELECT name from %s.sys.tables", dbName)
}
