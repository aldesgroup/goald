package goald

// Helps adapt to several types of SQL databases
type iDBAdapter interface {
	getTablesQuery(dbName string) string
}

func getAdapter(dbType databaseType) iDBAdapter {
	switch dbType {
	case dbTypeSQLSERVER:
		return &dbAdapterMSSQL{}
	default:
		panicf("Unhandled DB type: %s", dbType)
	}

	return nil
}
