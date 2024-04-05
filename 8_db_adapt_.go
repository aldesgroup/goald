package goald

import "github.com/aldesgroup/goald/features/utils"

// Helps adapt to several types of SQL databases
type iDBAdapter interface {
	getTablesQuery(dbName string) string
}

func getAdapter(dbType databaseType) iDBAdapter {
	switch dbType {
	case dbTypeSQLSERVER:
		return &dbAdapterMSSQL{}
	default:
		utils.Panicf("Unhandled DB type: %s", dbType)
	}

	return nil
}
