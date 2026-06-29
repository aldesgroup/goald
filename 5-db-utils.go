package goald

import "github.com/aldesgroup/goald/features/logging"

// FetchStringColumn executes a query that should only return an array of string (1 column)
func (thisDB *DB) FetchStringColumn(logger logging.ILogger, failIfErr bool, query string, args ...interface{}) (results []string) {
	// executing the query
	rows, errQuery := thisDB.Query(logger, query, args...)
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
