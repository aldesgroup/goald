package goald

import "github.com/aldesgroup/goald/features/logging"

// ----------------------------------------------------------------------------
// DB errors
// ----------------------------------------------------------------------------

type DbError int

const (
	DbErrorUNIDENTIFIED    DbError = 0
	DbErrorDUPLICATExENTRY DbError = 1
	DbErrorINVALIDxENTRY   DbError = 2
	DbErrorMISSINGxVALUE   DbError = 3
)

var dbErrors = map[int]string{
	int(DbErrorUNIDENTIFIED):    "unidentified error",
	int(DbErrorDUPLICATExENTRY): "duplicate entry error",
	int(DbErrorINVALIDxENTRY):   "invalid entry error",
	int(DbErrorMISSINGxVALUE):   "missing value error",
}

func (thisError DbError) String() string {
	return dbErrors[int(thisError)]
}

// Val helps implement the IEnum interface
func (thisError DbError) Val() int {
	return int(thisError)
}

// Values helps implement the IEnum interface
func (thisError DbError) Values() map[int]string {
	return dbErrors
}

// ----------------------------------------------------------------------------
// DB methods to quickly fetch stuff from the database
// ----------------------------------------------------------------------------

// FetchStringColumn executes a query that should only return an array of string (1 column)
func (thisDB *DB) FetchStringColumn(logger logging.ILogger, failIfErr bool, m mask, query string, args ...any) (results []string) {
	// executing the query
	rows, errQuery := thisDB.query(logger, nil, m, query, args...)
	if errQuery != nil {
		logger.Error(failIfErr, "Error while executing query '%s': %s", query, errQuery)
		return
	}

	// a temp result variable to hold the scanned value
	var result string

	// avoid forgetting to close the rows when exiting this function
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logger.Error(failIfErr, "Error while closing rows: %s", errClose)
		}
	}()

	// iterating over the result set
	for rows.Next() {
		if errScan := rows.Scan(&result); errScan != nil {
			logger.Error(failIfErr, "Error while scanning a row: %s", errScan)
			return
		}

		results = append(results, result)
	}

	// handling the error occurring during the call to .Next()
	if errNext := rows.Err(); errNext != nil {
		logger.Error(failIfErr, "Error while iterating over the rows: %s", errNext)
	}

	return
}

// FetchStringMap executes a query that should only return a map of string -> string
func (thisDB *DB) FetchStringMap(logger logging.ILogger, failIfErr bool, m mask, query string, args ...any) (results map[string]string) {
	// executing the query
	rows, errQuery := thisDB.query(logger, nil, m, query, args...)
	if errQuery != nil {
		logger.Error(failIfErr, "Error while executing query '%s': %s", query, errQuery)
		return
	}

	// preparing the returned map
	results = make(map[string]string)

	// a temp key->value pair to hold the scanned values
	var key, value string

	// avoid forgetting to close the rows when exiting this function
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logger.Error(failIfErr, "Error while closing rows: %s", errClose)
		}
	}()

	// iterating over the result set
	for rows.Next() { // iterating over the result set
		if errScan := rows.Scan(&key, &value); errScan != nil {
			logger.Error(failIfErr, "Error while scanning a row: %s", errScan)
			return
		}

		results[key] = value
	}

	// handling the error occurring during the call to .Next()
	if errNext := rows.Err(); errNext != nil {
		logger.Error(failIfErr, "Error while iterating over the rows: %s", errNext)
	}

	return
}
