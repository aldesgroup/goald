package goald

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

func openDB(conf *dbConfig) *sql.DB {
	// getting the password
	password := os.Getenv(conf.Password)

	// building the connection string
	var connStr string
	switch conf.DbType {
	case dbTypeSQLSERVER:
		connStr = fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s",
			conf.User, password, conf.DbPort, conf.DbName)
	default:
		panicf("Unhandled DB type '%s'", conf.DbType)
	}

	// Creating the DB object
	startDB := time.Now()
	db, errOpen := sql.Open(string(conf.DbType), connStr)
	if errOpen != nil {
		panicf("Error opening DB '%s': %s", conf.DbID, errOpen)
	}

	// Pinging
	if errPing := db.PingContext(context.Background()); errPing != nil {
		panicf("Issue while testing the '%s' DB: %s", conf.DbID, errPing)
	}

	fmt.Printf("Established connection to DB '%s' in %s!\n", conf.DbID, time.Since(startDB))

	return db
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
