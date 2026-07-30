package pgsql

import (
	"fmt"
	"strings"

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

// ----------------------------------------------------------------------------
// DB error parsing
// ----------------------------------------------------------------------------

func (thisAdapter *dbAdapterPGSQL) ParseDbError(err error) (goald.DbError, string) {

	// Duplicate entry error: `ERROR: duplicate key value violates unique constraint "uk__purchase_order__order_ref" (SQLSTATE 23505)``
	if strings.HasSuffix(err.Error(), "(SQLSTATE 23505)") {
		content := err.Error()
		start := strings.Index(content, "\"")
		end := strings.LastIndex(content, "\"")
		if start >= 0 && end > start {
			return goald.DbErrorDUPLICATExENTRY, content[start+1 : end]
		}

		return goald.DbErrorDUPLICATExENTRY, ""
	}

	// Invalid entry error: `ERROR: insert or update on table "purchase_order" violates foreign key constraint "fk__purchase_order__company__id" (SQLSTATE 23503)`
	if strings.HasSuffix(err.Error(), "(SQLSTATE 23503)") {
		content := err.Error()
		start1 := strings.Index(content, "\"")
		end1 := strings.Index(content[start1+1:], "\"") + start1 + 1
		start2 := strings.Index(content[end1+1:], "\"") + end1 + 1
		end2 := strings.LastIndex(content[start2+1:], "\"") + start2 + 1
		if start1 >= 0 && end1 > start1 && start2 >= 0 && end2 > start2 {
			return goald.DbErrorINVALIDxENTRY, content[start2+1 : end2]
		}
		return goald.DbErrorINVALIDxENTRY, ""
	}

	// Missing error: `ERROR: null value in column "main_contact__id" of relation "purchase_order" violates not-null constraint (SQLSTATE 23502)
	if strings.HasSuffix(err.Error(), "(SQLSTATE 23502)") {
		content := err.Error()
		start1 := strings.Index(content, "\"")
		end1 := strings.Index(content[start1+1:], "\"") + start1 + 1
		start2 := strings.Index(content[end1+1:], "\"") + end1 + 1
		end2 := strings.LastIndex(content[start2+1:], "\"") + start2 + 1
		if start1 >= 0 && end1 > start1 && start2 >= 0 && end2 > start2 {
			return goald.DbErrorMISSINGxVALUE, content[start2+1:end2] + "." + content[start1+1:end1]
		}

		return goald.DbErrorMISSINGxVALUE, ""
	}

	return goald.DbErrorUNIDENTIFIED, ""
}
