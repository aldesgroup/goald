package pgsql

import (
	"fmt"

	"github.com/aldesgroup/goald/features/dbconn"
)

// ----------------------------------------------------------------------------
// DB schema accesses
// ----------------------------------------------------------------------------

// GrantUsageCreateOnSchemaQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) GrantUsageCreateOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string {
	return fmt.Sprintf("GRANT USAGE, CREATE ON SCHEMA \"%s\" TO \"%s\"", schema, user)
}

// GrantUsageSchemaQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) GrantUsageOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string {
	return fmt.Sprintf("GRANT USAGE ON SCHEMA \"%s\" TO \"%s\"", schema, user)
}

// GrantAllPrivilegesOnSchemaQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) GrantAllPrivilegesOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string {
	return fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA \"%s\" TO \"%s\"", schema, user)
}

// GrantReadOnSchemaQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) GrantReadOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string {
	return fmt.Sprintf("GRANT SELECT, REFERENCES ON ALL TABLES IN SCHEMA \"%s\" TO \"%s\"", schema, user)
}

// ----------------------------------------------------------------------------
// DB schema discovery queries
// ----------------------------------------------------------------------------

// TablesQuery implements [goald.iDBAdapter]
func (thisAdapter *dbAdapterPGSQL) TablesQuery() string {
	return `
SELECT table_name
  FROM information_schema.tables
 WHERE table_schema = $1
   AND table_type = 'BASE TABLE'
 ORDER BY table_name;`
}

// ColumnsQuery implements [goald.iDBAdapter]
func (thisAdapter *dbAdapterPGSQL) ColumnsQuery() string {
	return `
SELECT table_name, column_name, is_nullable, character_maximum_length, numeric_precision, numeric_scale, datetime_precision, data_type
  FROM information_schema.columns
 WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
   AND table_schema = $1;`
}

// ForeignKeysQuery implements [goald.iDBAdapter]
func (thisAdapter *dbAdapterPGSQL) ForeignKeysQuery(fkPrefix string) string {
	return fmt.Sprintf(`
SELECT rc.constraint_name, tc.table_schema || '.' || tc.table_name
  FROM information_schema.referential_constraints rc
  JOIN information_schema.table_constraints tc
    ON tc.constraint_name = rc.constraint_name
   AND tc.constraint_schema = rc.constraint_schema
 WHERE rc.constraint_schema = $1
   AND rc.constraint_name LIKE '%s%%'`, fkPrefix)
}

// UniqueConstraintsQuery(ukPrefix string) string implements [goald.iDBAdapter]
func (thisAdapter *dbAdapterPGSQL) UniqueConstraintsQuery(prefix string) string {
	return fmt.Sprintf(`
SELECT constraint_name, table_name 
  FROM information_schema.table_constraints
 WHERE constraint_type = 'UNIQUE' 
   AND constraint_name LIKE '%s%%' 
   AND constraint_schema = $1`, prefix)
}
