// ------------------------------------------------------------------------------------------------
// Here is the global code generation routine
// ------------------------------------------------------------------------------------------------
package goald

import (
	"log"
	"os"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

type codeGenLevel int

const codeGenOBJECTS codeGenLevel = 1
const codeGenCLASSES codeGenLevel = 2

// this function shows that our server, when run in dev with the right arguments,
// can be used as a development server, generating code for us
func (thisServer *server) runCodeGen(srcdir string, level codeGenLevel, _ bool, regen bool) {
	switch level {
	case codeGenOBJECTS:
		start := time.Now()

		// TODO optimize with go routines here (?)

		// we're making all the databases globally accessible
		thisServer.generateDatabasesList(srcdir)

		// we're making all the business objects reachable with the `reflect` package this way
		thisServer.generateObjectRegistry(srcdir, ".", false, map[string]*businessObjectEntry{}, regen)

		log.Printf("done generating the DB & BO registries in %s", time.Since(start))
		os.Exit(0)

	case codeGenCLASSES:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build a static representation of our BOs
		thisServer.generateObjectClasses(srcdir, regen)

		// now, using the `reflect` package, we can "easily" build utils for our BOs,
		// that should help us avoid using the `reflect` package at runtime;
		thisServer.generateObjectUtils(srcdir, ".", regen)

		// one reason for having "classes", is to be able to instantiate them, using constructors;
		// we can do this right away, no need for a 3rd building step
		// generateObjectConstructors(srcdir, currentPath, map[string]bool{})

		log.Printf("done generating the classes & constructors in %s", time.Since(start))
		os.Exit(0)

	default:
		utils.Panicf("Not handling to code generation level: %d", level)
	}
}
