package goald

import "github.com/aldesgroup/goald/features/logging"

// ----------------------------------------------------------------------------
// DB errors
// ----------------------------------------------------------------------------

// DbError is a string-based error type, so that its values can be constants, comparable in a switch or with errors.Is
type DbError string

// Error implements the error interface
func (thisDbError DbError) Error() string {
	return string(thisDbError)
}

const (
	DbErrorUNIDENTIFIED    = DbError("unidentified error")
	DbErrorDUPLICATExENTRY = DbError("duplicate entry error")
	DbErrorINVALIDxENTRY   = DbError("invalid entry error")
	DbErrorMISSINGxVALUE   = DbError("missing value error")
)

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
