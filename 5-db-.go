package goald

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/logging"
	_ "github.com/microsoft/go-mssqldb"
)

// ------------------------------------------------------------------------------------------------
// Listing the DB types Goald can handle
// ------------------------------------------------------------------------------------------------

const (
	// DbTypeSQLSERVER  = "sqlserver"
	DbTypePOSTGRESQL dbconn.DatabaseType = "postgresql"
)

var allDbTypes = []dbconn.DatabaseType{
	// DbTypeSQLSERVER,
	DbTypePOSTGRESQL,
}

// ------------------------------------------------------------------------------------------------
// Describing how a DB adapter should behave
// ------------------------------------------------------------------------------------------------

// Helps adapt to several types of SQL databases
type iDBAdapter interface {
	// DB server-related methods
	DatabaseType() dbconn.DatabaseType                                                      // the type of DB this adapter is for - should match what's configured in aldev config
	DriverName() string                                                                     // the name of the driver to use for this DB type
	ConnectionString(dbConfig *dbconn.DbConfig, user dbconn.DbUserName, pass string) string // building the connection string for a given DB config and user/pass

	// DB server-related queries
	SchemaExistsQuery() string                                  // checking if a given schema exists in the DB server
	UserExistsQuery() string                                    // checking if a given user exists in the DB server
	CreateUserQuery(user dbconn.DbUserName, pass string) string // creating a user in the DB server

	// DB schema-related authorization queries
	GrantUsageCreateOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string   // granting privileges to a user on a schema
	GrantUsageOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string         // granting privileges to a user on a schema
	GrantAllPrivilegesOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string // granting all privileges to a user on all tables in a schema
	GrantReadOnSchemaQuery(schema dbconn.DbSchemaName, user dbconn.DbUserName) string          // granting read privileges to a user on all tables in a schema

	// DB schema-related discovery queries
	TablesQuery() string                           // retrieving the existing table names in a given schema
	ColumnsQuery() string                          // retrieving columns info in a given schema
	ForeignKeysQuery(fkPrefix string) string       // retrieving foreign keys info in a given schema
	UniqueConstraintsQuery(ukPrefix string) string // retrieving unique constraints info in a given schema

	// DB table-related modification queries
	DropTableFkQuery(tableName string, fkName string) string                         // dropping the foreign keys of a given table
	AddNotNullQuery(tableName string, columnName string, columnType string) string   // adding a NOT NULL constraint to a given column of a given table
	DropNotNullQuery(tableName string, columnName string, columnType string) string  // dropping a NOT NULL constraint from a given column of a given table
	ModifyColumnQuery(tableName string, columnName string, columnType string) string // modifying a given column of a given table

	// DB column-related building methods
	SQLColumnDeclaration(property IBusinessObjectProperty) (string, string) // building the SQL column declaration for a given BO property
}

// ------------------------------------------------------------------------------------------------
// Goald databases
// ------------------------------------------------------------------------------------------------

// goald's own DB object
type DB struct {
	*sql.DB
	get    iDBAdapter
	config *dbconn.DbSchemaConfig
}

// accessing the config
func (thisDB *DB) GetConfig() *dbconn.DbSchemaConfig {
	return thisDB.config
}

// ------------------------------------------------------------------------------------------------
// Low level DB operations
// ------------------------------------------------------------------------------------------------

func logSQL(logger logging.ILogger, db *DB, start time.Time, query string, args ...any) {
	if logger.IsVerbose() {
		// TODO later: plug in a mechanism of gathering analytics here
		logger.Debug(fmt.Sprintf("Run from '%s' in %s (with args: %+v): %s", db.config.Name, time.Since(start), args, strings.ReplaceAll(query, "\n", "\n-   ")))
	}
}

// proxying this function so as to add functionality
func (thisDB *DB) Query(logger logging.ILogger, query string, args ...any) (*sql.Rows, error) {
	defer logSQL(logger, thisDB, time.Now(), query, args...)
	return thisDB.DB.Query(query, args...)
}

// proxying this function so as to add functionality
func (thisDB *DB) Exec(logger logging.ILogger, query string, args ...any) (sql.Result, error) {
	defer logSQL(logger, thisDB, time.Now(), query, args...)
	return thisDB.DB.Exec(query, args...)
}

// shortcut for executing a query and panicking if it fails
func (thisDB *DB) MustQuery(logger logging.ILogger, query string, args ...any) *sql.Rows {
	rows, err := thisDB.Query(logger, query, args...)
	core.PanicMsgIfErr(err, "Error executing SQL statement: '%s' with args: %+v", query, args)
	return rows
}

// shortcut for executing a query and panicking if it fails
func (thisDB *DB) MustExec(logger logging.ILogger, query string, args ...any) sql.Result {
	result, err := thisDB.Exec(logger, query, args...)
	core.PanicMsgIfErr(err, "Error executing SQL statement: '%s' with args: %+v", query, args)
	return result
}

// ------------------------------------------------------------------------------------------------
// Opening a DB, checking it, etc.
// ------------------------------------------------------------------------------------------------

func (thisServer *server) connectDbSchema(dbSchema *dbconn.DbSchemaConfig) {
	// LFG
	start := time.Now()

	// getting the right adapter for the current DB config
	adapter := getDbAdapter(dbSchema.DbConfig.Type)

	// building the connection string
	connStr := adapter.ConnectionString(dbSchema.DbConfig, dbSchema.User, dbSchema.Pass)

	// opening the DB connection
	db, err := sql.Open(adapter.DriverName(), connStr)
	core.PanicMsgIfErr(err, "Error opening DB schema '%s' from DB '%s'", dbSchema.Name, dbSchema.DbConfig.Database)

	// checking the connection to the DB schema
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	core.PanicMsgIfErr(db.PingContext(ctx), "Error pinging DB schema '%s' from DB '%s'", dbSchema.Name, dbSchema.DbConfig.Database)

	// registration for later use
	goaldDB := GetDB(dbSchema.Name)
	goaldDB.get = adapter
	goaldDB.config = dbSchema
	goaldDB.DB = db

	// bit of logging
	thisServer.Info(fmt.Sprintf("Established connection to DB schema '%s' in %s", dbSchema.Name, time.Since(start)))
}
