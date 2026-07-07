package dbconn

// ------------------------------------------------------------------------------------------------
// Describing DB connections
// ------------------------------------------------------------------------------------------------

type DatabaseType string
type DbSchemaName string
type DbUserName string

type DbConfig struct {
	Type     DatabaseType
	Host     string
	Port     int
	Database string
	SSLMode  string // e.g. "require" for Azure, empty = pgx default (prefer)
	Admin    *struct {
		User DbUserName
		Pass string
	}
	Schemas map[DbSchemaName]*DbSchemaConfig
}

type schemaAccessType string

const (
	SchemaAccessREAD  schemaAccessType = "read"
	SchemaAccessWRITE schemaAccessType = "write"
)

type DbSchemaConfig struct {
	Name     DbSchemaName
	User     DbUserName
	Pass     string
	Access   map[DbUserName]schemaAccessType // access to this schema for other users (user -> access type)
	DbConfig *DbConfig                       // back reference to the parent DB config
}
