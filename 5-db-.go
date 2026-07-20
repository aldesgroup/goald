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

	SupportsReturningID() bool // tells whether INSERT statements can use a 'RETURNING' clause to get the new row's ID

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
	name   dbconn.DbSchemaName
	do     *sql.DB
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

// a mask, when not nil, allows to hide some of the arguments of a query from the logs - useful for sensitive data like passwords
type mask []bool

func logSQL(logger logging.ILogger, db *DB, start time.Time, m mask, query string, args ...any) {
	if logger.IsVerbose() {
		// TODO later: plug in a mechanism of gathering analytics here
		var loggedArgs []any
		for i, arg := range args {
			if m != nil && i < len(m) && m[i] {
				loggedArgs = append(loggedArgs, "*****")
			} else {
				loggedArgs = append(loggedArgs, arg)
			}
		}
		logger.Debug(fmt.Sprintf("Run from '%s' in %s (with args: %+v): %s", db.config.Name,
			time.Since(start), loggedArgs, strings.Join(strings.Fields(query), " ")))
	}
}

// proxying this function so as to add functionality
func (thisDB *DB) query(logger logging.ILogger, m mask, query string, args ...any) (*sql.Rows, error) {
	defer logSQL(logger, thisDB, time.Now(), m, query, args...)
	return thisDB.do.Query(query, args...)
}

// proxying this function so as to add functionality
func (thisDB *DB) exec(logger logging.ILogger, m mask, query string, args ...any) (sql.Result, error) {
	defer logSQL(logger, thisDB, time.Now(), m, query, args...)
	return thisDB.do.Exec(query, args...)
}

// proxying this function so as to add functionality - used for the 'RETURNING id'-style insert queries
func (thisDB *DB) queryRow(logger logging.ILogger, m mask, query string, args ...any) *sql.Row {
	defer logSQL(logger, thisDB, time.Now(), m, query, args...)
	return thisDB.do.QueryRow(query, args...)
}

// shortcut for executing a query and panicking if it fails
func (thisDB *DB) mustQuery(logger logging.ILogger, m mask, query string, args ...any) *sql.Rows {
	rows, err := thisDB.query(logger, m, query, args...)
	core.PanicMsgIfErr(err, "Error executing SQL statement: '%s' with args: %+v", query, args)
	return rows
}

// shortcut for executing a query and panicking if it fails
func (thisDB *DB) mustExec(logger logging.ILogger, m mask, query string, args ...any) sql.Result {
	result, err := thisDB.exec(logger, m, query, args...)
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
	goaldDB.do = db
	goaldDB.get = adapter
	goaldDB.config = dbSchema

	// bit of logging
	thisServer.Info(fmt.Sprintf("Established connection to DB schema '%s' in %s", dbSchema.Name, time.Since(start)))
}

// ------------------------------------------------------------------------------------------------
// Misc
// ------------------------------------------------------------------------------------------------

func DbFor(bObj IBusinessObject) *DB {
	return modelForName(bObj.GetClassName(bObj)).getInDB()
}
