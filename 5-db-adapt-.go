package goald

// Helps adapt to several types of SQL databases
type iDBAdapter interface {
	getConnectionString(conf *dbConfig) string
	getTablesQuery(dbName string) string
	getSQLColumnDeclaration(property iBusinessObjectProperty) string
}
