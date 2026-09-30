package goald

import (
	"database/sql"
	"fmt"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
)

func (thisServer *server) initDbServer(dbConfig *dbconn.DbServerConfig) {
	// getting the DB adapter for this DB type
	adapter := getDbAdapter(dbConfig.Type)

	// building the connection string
	connStr := adapter.ConnectionString(dbConfig, dbConfig.Admin.User, dbConfig.Admin.Pass)

	// opening the DB connection
	db, err := sql.Open(adapter.DriverName(), connStr)
	core.PanicMsgIfErr(err, "Error opening DB connection")

	// testing the DB connection - retrying this 10 times with a 1 second delay between each retry
	core.PanicMsgIfErr(db.Ping(), "Error pinging DB")
	thisServer.Info("Successfully connected to " + string(dbConfig.Type) + " DB server: " + dbConfig.Database)

	// the function we'll use to execute SQL statements
	execSQL := func(sqlStmt string, args ...any) {
		query := fmt.Sprintf(sqlStmt, args...)
		_, err := db.Exec(fmt.Sprintf(sqlStmt, args...))
		core.PanicMsgIfErr(err, "Error executing SQL statement: '%s'", query)
		thisServer.Debug("Successfully run: " + query)
	}

	// creating the schemas and users for each schema defined in the config
	for _, schema := range dbConfig.Schemas {
		// creating the schema only if it doesn't exist yet
		var schemaExists bool
		err = db.QueryRow(adapter.SchemaExistsQuery(), schema.Name).Scan(&schemaExists)
		core.PanicMsgIfErr(err, "Error checking if schema '%s' exists", schema.Name)
		if !schemaExists {
			execSQL("CREATE SCHEMA \"%s\"", schema.Name)
		}

		// creating the user for the schema only if it doesn't exist yet
		var userExists bool
		err = db.QueryRow(adapter.UserExistsQuery(), schema.User).Scan(&userExists)
		core.PanicMsgIfErr(err, "Error checking if user '%s' exists", schema.User)
		if !userExists {
			// creating the user for the schema
			execSQL("%s", adapter.CreateUserQuery(schema.User, schema.Pass))

			// granting privileges to the user on the schema
			execSQL("%s", adapter.GrantUsageCreateOnSchemaQuery(schema.Name, schema.User))
		}
	}

	// granting access between schemas if specified
	for _, schema := range dbConfig.Schemas {
		for otherUser := range schema.Access {
			// granting privileges to the user on the other schema
			execSQL("%s", adapter.GrantUsageOnSchemaQuery(schema.Name, otherUser))
		}
	}
}
