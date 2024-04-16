// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the DB list
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"strings"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

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

const dbFOLDER = "_generated/db"
const dbFILE = "db-list.go"

func (thisServer *server) generateDatabasesList(srcdir string) {
	start := time.Now()

	// starting to build the file content, with the same context
	content := `package db`

	if len(thisServer.config.commonPart().Databases) > 0 {
		content += `

import (
	"sync"

	g "github.com/aldesgroup/goald"
)

`

		for _, dbConfig := range thisServer.config.commonPart().Databases {
			dbParagraph := strings.ReplaceAll(dbTEMPLATE, "$$dbID$$", utils.PascalToCamel(string(dbConfig.DbID)))
			dbParagraph = strings.ReplaceAll(dbParagraph, "$$DbID$$", utils.ToPascal(string(dbConfig.DbID)))
			dbParagraph = strings.ReplaceAll(dbParagraph, "$$realDbID$$", string(dbConfig.DbID))
			content += dbParagraph + newline
		}

		// writing to file
		utils.WriteToFile(content, srcdir, dbFOLDER, dbFILE)
		println(fmt.Sprintf("DB list generated in %s", time.Since(start)))
	}
}
