// ------------------------------------------------------------------------------------------------
// Here is the code used for generating the DB list
// ------------------------------------------------------------------------------------------------
package goald

import (
	"log"
	"os"
	"path"
	"strings"
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

const dbFOLDER = "db"
const dbFILE = "db_list.go"

func (thisServer *server) generateDatabasesList(srcdir string) {
	// checking the class folder exist, or creating it on the way
	dbDir := path.Join(srcdir, dbFOLDER)
	if !DirExists(dbDir) {
		panicErrf(os.Mkdir(dbDir, 0o777), "Could not create the db folder '%s'", dbDir)
	}

	// creating the file
	fileName := path.Join(dbDir, dbFILE)

	file, errCreate := os.Create(fileName)
	if errCreate != nil {
		panicf("Could not create file %s; cause: %s", fileName, errCreate)
	}

	defer func() {
		if errClose := file.Close(); errClose != nil {
			log.Fatalf("Could not properly close file %s; cause: %s", fileName, errClose)
		}
	}()

	// starting to build the file content, with the same context
	content := `package db

import (
	"sync"

	g "git-ext.aldes.com/j.wan/arch-poc/goald"
)

`

	for _, dbConfig := range thisServer.config.getCommonConfig().Databases {
		dbParagraph := strings.ReplaceAll(dbTEMPLATE, "$$dbID$$", PascalToCamel(string(dbConfig.DbID)))
		dbParagraph = strings.ReplaceAll(dbParagraph, "$$DbID$$", ToPascal(string(dbConfig.DbID)))
		dbParagraph = strings.ReplaceAll(dbParagraph, "$$realDbID$$", string(dbConfig.DbID))
		content += dbParagraph + newline
	}

	// writing to file
	if _, errWrite := file.WriteString(content); errWrite != nil {
		panicErrf(errWrite, "Could not write file '%s'", fileName)
	}
}
