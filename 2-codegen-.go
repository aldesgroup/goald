// ------------------------------------------------------------------------------------------------
// Here is the global code generation routine
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

type codeGenLevel int
type packageName string

const codeGenOBJECTS codeGenLevel = 1
const codeGenCLASSES codeGenLevel = 2
const codeGenUTILS codeGenLevel = 3

// this function shows that our server, when run in dev with the right arguments,
// can be used as a development server, generating code for us
func (thisServer *server) runCodeGen(srcdir string, level codeGenLevel, webdir string, regen bool) {
	switch level {
	case codeGenOBJECTS:
		start := time.Now()

		// TODO optimize with go routines here (?)

		// we're making all the databases globally accessible
		thisServer.generateDatabasesList(srcdir)

		// generating the class utils and the packages that register them, and make the corresponding business objects "importable"
		thisServer.generateIncludes(srcdir, ".", false, map[packageName]map[className]*classUtilsCore{}, regen)

		slog.Info(fmt.Sprintf("done generating the DB & BO registries in %s", time.Since(start)))
		os.Exit(0)

	case codeGenCLASSES:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build a static representation of our BOs
		thisServer.generateObjectClasses(srcdir, regen)

		slog.Info(fmt.Sprintf("done generating the BO classes in %s", time.Since(start)))
		os.Exit(0)

	case codeGenUTILS:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build utils for our BOs,
		// that should help us avoid using the `reflect` package at runtime;
		thisServer.generateObjectValueMappers(srcdir, ".", regen)
		thisServer.generateWebAppModels(webdir, regen)

		slog.Info(fmt.Sprintf("done generating the BO utils in %s", time.Since(start)))
		os.Exit(0)

	default:
		utils.Panicf("Not handling to code generation level: %d", level)
	}
}

// ------------------------------------------------------------------------------------------------
// Utilities - this is the only place we're allowed to use the 'reflect' package
// ------------------------------------------------------------------------------------------------
