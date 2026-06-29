package goald

import (
	"context"
	"database/sql"
	"fmt"
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
	// Registration
	GetDatabaseType() dbconn.DatabaseType

	// DB server init
	InitDbServer(logger logging.ILogger, dbConfig *dbconn.DbConfig)

	// DB schema opening
	OpenDbSchema(logger logging.ILogger, dbSchema *dbconn.DbSchemaConfig) (*sql.DB, error)

	// Specific queries
	GetTablesQuery() string
}

// ------------------------------------------------------------------------------------------------
// Goald databases
// ------------------------------------------------------------------------------------------------

// goald's own DB object
type DB struct {
	*sql.DB
	iDBAdapter
	name   dbconn.DbSchemaName
	config *dbconn.DbSchemaConfig
}

func logSQL(logger logging.ILogger, start time.Time, query string, args ...any) {
	if logger.IsVerbose() {
		// TODO later: plug in a mechanism of gathering analytics here
		logger.Debug(fmt.Sprintf("Run in %s: %s (with args: %+v)", time.Since(start), query, args))
	}
}

// proxying this function so as to add functionality
func (thisDB *DB) Query(logger logging.ILogger, query string, args ...any) (*sql.Rows, error) {
	defer logSQL(logger, time.Now(), query, args...)
	return thisDB.DB.Query(query, args...)
}

// proxying this function so as to add functionality
func (thisDB *DB) Exec(logger logging.ILogger, query string, args ...any) (sql.Result, error) {
	defer logSQL(logger, time.Now(), query, args...)
	return thisDB.DB.Exec(query, args...)
}

// ------------------------------------------------------------------------------------------------
// Opening a DB, checking it, etc.
// ------------------------------------------------------------------------------------------------

func (thisServer *server) connectDbSchema(dbSchema *dbconn.DbSchemaConfig) {
	// LFG
	start := time.Now()

	// getting the right adapter for the current DB config
	adapter := getDbAdapter(dbSchema.DbConfig.Type)

	// opening the DB schema
	db, err := adapter.OpenDbSchema(thisServer, dbSchema)
	core.PanicMsgIfErr(err, "Error opening DB schema '%s' from DB '%s'", dbSchema.Name, dbSchema.DbConfig.Database)

	// checking the connection to the DB schema
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	core.PanicMsgIfErr(db.PingContext(ctx), "Error pinging DB schema '%s' from DB '%s'", dbSchema.Name, dbSchema.DbConfig.Database)

	// registration for later use
	goaldDB := GetDB(dbSchema.Name)
	goaldDB.iDBAdapter = adapter
	goaldDB.config = dbSchema
	goaldDB.DB = db

	// bit of logging
	thisServer.Info(fmt.Sprintf("Established connection to DB schema '%s' in %s", dbSchema.Name, time.Since(start)))
}

// // ------------------------------------------------------------------------------------------------
// // Quick DB operations
// // ------------------------------------------------------------------------------------------------

// // Executes a query that should only return an array of string (1 column)
// func (thisDB *DB) FetchStringColumn(query string, args ...interface{}) (results []string, err error) {
// 	// TODO better handle logging
// 	rows, err := thisDB.Query(query, args...)
// 	if err != nil {
// 		return nil, ErrorC(err, "Error while executing query '%s': %s", query, err)
// 	}

// 	// making sure we're closing the rows
// 	defer func() {

// 		if errClose := rows.Close(); errClose != nil {
// 			// TODO do something
// 			println(errClose)
// 		}
// 	}()

// 	// iterating over the result set
// 	var result string
// 	for rows.Next() {
// 		if err = rows.Scan(&result); err != nil {
// 			return nil, ErrorC(err, "Error while scanning a row: %s", err)
// 		}

// 		results = append(results, result)
// 	}

// 	// handling the error occurring during the call to .Next()
// 	if err = rows.Err(); err != nil {
// 		return nil, ErrorC(err, "Error while iterating over the rows: %s", err)
// 	}

// 	return
// }
