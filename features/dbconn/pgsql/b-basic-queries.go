package pgsql

// getTablesQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) GetTablesQuery() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema = $1 AND table_type = 'BASE TABLE' ORDER BY table_name;"
}
