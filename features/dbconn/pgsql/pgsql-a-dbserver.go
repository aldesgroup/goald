package pgsql

import (
	"fmt"

	"github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/dbconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ----------------------------------------------------------------------------
// DB Adapter declaration and registration
// ----------------------------------------------------------------------------

type dbAdapterPGSQL struct{}

func init() {
	goald.RegisterDbAdapter(&dbAdapterPGSQL{})
}

// DatabaseType implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) DatabaseType() dbconn.DatabaseType {
	return goald.DbTypePOSTGRESQL
}

// DriverName implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) DriverName() string {
	return "pgx"
}

// SupportsReturningID implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) SupportsReturningID() bool {
	return true
}

// ConnectionString implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) ConnectionString(dbConfig *dbconn.DbServerConfig, user dbconn.DbUserName, pass string) string {
	sslMode := dbConfig.SSLMode
	if sslMode == "" {
		sslMode = "prefer"
	}
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		pass,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		sslMode,
	)
}

// ----------------------------------------------------------------------------
// DB server queries
// ----------------------------------------------------------------------------

// SchemaExistsQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) SchemaExistsQuery() string {
	return "SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)"
}

// UserExistsQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) UserExistsQuery() string {
	return "SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)"
}

// CreateUserQuery implements [goald.iDBAdapter].
func (thisAdapter *dbAdapterPGSQL) CreateUserQuery(user dbconn.DbUserName, pass string) string {
	return fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s'", user, pass)
}
