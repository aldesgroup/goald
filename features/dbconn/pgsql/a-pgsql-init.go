package pgsql

import (
	"database/sql"
	"fmt"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/logging"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ----------------------------------------------------------------------------
// DB Adapter declaration and registration
// ----------------------------------------------------------------------------

type dbAdapterPGSQL struct{}

func (thisAdapter *dbAdapterPGSQL) GetDatabaseType() dbconn.DatabaseType {
	return goald.DbTypePOSTGRESQL
}

func init() {
	goald.RegisterDbAdapter(&dbAdapterPGSQL{})
}

// ----------------------------------------------------------------------------
// DB server init
// ----------------------------------------------------------------------------

func (thisAdapter *dbAdapterPGSQL) InitDbServer(logger logging.ILogger, dbConfig *dbconn.DbConfig) {
	// building the connection string
	sslMode := dbConfig.SSLMode
	if sslMode == "" {
		sslMode = "prefer"
	}
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.Admin.User,
		dbConfig.Admin.Pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		sslMode,
	)

	// opening the DB connection
	db, err := sql.Open("pgx", connStr)
	core.PanicMsgIfErr(err, "Error opening DB connection")

	// testing the DB connection - retrying this 10 times with a 1 second delay between each retry
	core.PanicMsgIfErr(db.Ping(), "Error pinging DB")
	logger.Info("Successfully connected to the DB")

	// the function we'll use to execute SQL statements
	execSQL := func(sqlStmt string, args ...any) {
		query := fmt.Sprintf(sqlStmt, args...)
		_, err := db.Exec(fmt.Sprintf(sqlStmt, args...))
		core.PanicMsgIfErr(err, "Error executing SQL statement: '%s'", query)
		logger.Info("Successfully run: " + query)
	}

	// leeping track of the newly created users
	newUsers := map[string]bool{}

	// creating the schemas and users for each schema defined in the config
	for _, schema := range dbConfig.Schemas {
		// creating the schema only if it doesn't exist yet
		var schemaExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)", schema.Name).Scan(&schemaExists)
		core.PanicMsgIfErr(err, "Error checking if schema '%s' exists", schema.Name)
		if !schemaExists {
			execSQL("CREATE SCHEMA \"%s\"", schema.Name)
		}

		// creating the user for the schema only if it doesn't exist yet
		var userExists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)", schema.User).Scan(&userExists)
		core.PanicMsgIfErr(err, "Error checking if user '%s' exists", schema.User)
		if !userExists {
			// creating the user for the schema
			execSQL("CREATE USER \"%s\" WITH PASSWORD '%s'", schema.User, schema.Pass)
			newUsers[schema.User] = true

			// granting privileges to the user on the schema
			execSQL("GRANT USAGE, CREATE ON SCHEMA \"%s\" TO \"%s\"", schema.Name, schema.User)

			// granting privileges to the user on all tables in the schema
			execSQL("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA \"%s\" TO \"%s\"", schema.Name, schema.User)

		}
	}

	// granting access between schemas if specified
	for _, schema := range dbConfig.Schemas {
		for otherSchema, access := range schema.Access {
			if newUsers[schema.User] {
				execSQL("GRANT USAGE ON SCHEMA \"%s\" TO \"%s\"", otherSchema, schema.User)
				switch access {
				case dbconn.SchemaAccessWRITE:
					execSQL("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA \"%s\" TO \"%s\"", otherSchema, schema.User)
				case dbconn.SchemaAccessREAD:
					execSQL("GRANT SELECT ON ALL TABLES IN SCHEMA \"%s\" TO \"%s\"", otherSchema, schema.User)
				default:
					panic(fmt.Sprintf("Unknown access type '%s' for schema '%s'", access, otherSchema))
				}
			}
		}
	}
}

// ----------------------------------------------------------------------------
// DB schema connection
// ----------------------------------------------------------------------------

// this is just about opening the connection; pinging it is done in the calling function
func (thisAdapter *dbAdapterPGSQL) OpenDbSchema(logger logging.ILogger, dbSchema *dbconn.DbSchemaConfig) (*sql.DB, error) {
	// building the connection string
	sslMode := dbSchema.DbConfig.SSLMode
	if sslMode == "" {
		sslMode = "prefer"
	}
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		dbSchema.User,
		dbSchema.Pass,
		dbSchema.DbConfig.Host,
		dbSchema.DbConfig.Port,
		dbSchema.DbConfig.Database,
		sslMode,
	)

	// opening the DB connection
	return sql.Open("pgx", connStr)
}
