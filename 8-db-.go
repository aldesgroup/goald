package goald

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/aldesgroup/goald/features/utils"
	_ "github.com/microsoft/go-mssqldb"
)

// ------------------------------------------------------------------------------------------------
// Listing the known DB drivers
// ------------------------------------------------------------------------------------------------

type databaseType string

const (
	dbTypeSQLSERVER = "sqlserver"
)

// ------------------------------------------------------------------------------------------------
// Goald databases
// ------------------------------------------------------------------------------------------------

// goald's own DB object
type DB struct {
	*sql.DB
	config  *dbConfig
	adapter iDBAdapter
}

func logSQL(start time.Time, query string, args ...any) {
	log.Printf("Run in %s: %s (with args: %+v)", time.Since(start), query, args)
}

// proxying this function so as to add functionality
func (thisDB *DB) Query(query string, args ...any) (*sql.Rows, error) {
	defer logSQL(time.Now(), query, args...)
	return thisDB.DB.Query(query, args...)
}

// proxying this function so as to add functionality
func (thisDB *DB) Exec(query string, args ...any) (sql.Result, error) {
	defer logSQL(time.Now(), query, args...)
	return thisDB.DB.Exec(query, args...)
}

// ------------------------------------------------------------------------------------------------
// Opening a DB, checking it, etc.
// ------------------------------------------------------------------------------------------------

func openDB(conf *dbConfig) (*sql.DB, iDBAdapter) {
	// getting the right adapter for the current DB config
	var adapter iDBAdapter
	switch conf.DbType {
	case dbTypeSQLSERVER:
		adapter = &dbAdapterMSSQL{}
	default:
		utils.Panicf("Unhandled DB type: %s", conf.DbType)
	}

	// // making sure the DB exists
	// if conf.MakeExist {

	// }

	// connection string
	connStr := adapter.getConnectionString(conf)

	// Creating the DB object by opening connections with it
	startDB := time.Now()
	db, errOpen := sql.Open(string(conf.DbType), connStr)
	if errOpen != nil {
		utils.Panicf("Error opening DB '%s': %s", conf.DbID, errOpen)
	}

	// Pinging
	if errPing := db.Ping(); errPing != nil {
		utils.Panicf("Issue while testing the '%s' DB: %s", conf.DbID, errPing)
	}

	slog.Info(fmt.Sprintf("Established connection to DB '%s' in %s!\n", conf.DbID, time.Since(startDB)))

	return db, adapter
}

// ------------------------------------------------------------------------------------------------
// Quick DB operations
// ------------------------------------------------------------------------------------------------

// Executes a query that should only return an array of string (1 column)
func (thisDB *DB) FetchStringColumn(query string, args ...interface{}) (results []string, err error) {
	// TODO better handle logging
	rows, err := thisDB.Query(query, args...)
	if err != nil {
		return nil, ErrorC(err, "Error while executing query '%s': %s", query, err)
	}

	// making sure we're closing the rows
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			// TODO do something
			println(errClose)
		}
	}()

	// iterating over the result set
	var result string
	for rows.Next() {
		if err = rows.Scan(&result); err != nil {
			return nil, ErrorC(err, "Error while scanning a row: %s", err)
		}

		results = append(results, result)
	}

	// handling the error occurring during the call to .Next()
	if err = rows.Err(); err != nil {
		return nil, ErrorC(err, "Error while iterating over the rows: %s", err)
	}

	return
}
