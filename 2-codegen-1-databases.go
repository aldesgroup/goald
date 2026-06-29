// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the DB list
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
)

const dbFILExINIT = `

import (
	"sync"

	g "github.com/aldesgroup/goald"$$otherImports$$
)

`

const dbTEMPLATE = `// Access to the configured "$$realDbID$$" database

var $$dbID$$DB *g.DB
var $$dbID$$DBOnce sync.Once

func $$DbID$$() *g.DB {
	$$dbID$$DBOnce.Do(func() {
		$$dbID$$DB = g.GetDB("$$realDbID$$")
	})

	return $$dbID$$DB
}
`

const dbFOLDER = "_include/db"
const dbFILE = "db-list.go"

func (thisServer *server) generateDatabasesList(srcdir string) {
	start := time.Now()

	// starting to build the file content, with the same context
	content := `package db`

	// preparing for additional imports
	otherImportsMap := map[string]bool{}

	// adding 1 DB instance per DB schema configured in the server config, if any
	if len(thisServer.config.base().DBServers) > 0 {
		content += dbFILExINIT

		for _, dbConfig := range core.GetSortedValues(thisServer.config.base().DBServers) {
			otherImportsMap[getPackageForDbServerType(dbConfig.Type)] = true
			for _, schemaConfig := range core.GetSortedValues(dbConfig.Schemas) {
				dbParagraph := strings.ReplaceAll(dbTEMPLATE, "$$dbID$$", core.PascalToCamel(string(schemaConfig.Name)))
				dbParagraph = strings.ReplaceAll(dbParagraph, "$$DbID$$", core.ToPascal(string(schemaConfig.Name)))
				dbParagraph = strings.ReplaceAll(dbParagraph, "$$realDbID$$", string(schemaConfig.Name))
				content += dbParagraph + newline
			}
		}
	}

	// adding the necessary imports for the DB packages used
	otherImports := ""
	for _, dbPackage := range core.GetSortedKeys(otherImportsMap) {
		otherImports += fmt.Sprintf("\n\t_ \"github.com/aldesgroup/goald/features/dbconn/%s\"", dbPackage)
	}
	content = strings.ReplaceAll(content, "$$otherImports$$", otherImports)

	// writing to file
	core.WriteToFile(content, srcdir, dbFOLDER, dbFILE)
	println(fmt.Sprintf("DB list generated in %s", time.Since(start)))
}

func getPackageForDbServerType(databaseType dbconn.DatabaseType) string {
	switch databaseType {
	case "postgresql":
		return "pgsql"
	default:
		return ""
	}
}
