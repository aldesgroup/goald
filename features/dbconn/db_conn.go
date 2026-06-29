package dbconn

// ------------------------------------------------------------------------------------------------
// Describing DB connections
// ------------------------------------------------------------------------------------------------

type DatabaseType string

type DbSchemaName string

type DbConfig struct {
	Type     DatabaseType
	Host     string
	Port     int
	Database string
	SSLMode  string // e.g. "require" for Azure, empty = pgx default (prefer)
	Admin    *struct {
		User string
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
	User     string
	Pass     string
	Access   map[DbSchemaName]schemaAccessType // access to other schemas, to allow for easy joins
	DbConfig *DbConfig                         // back reference to the parent DB config
}
