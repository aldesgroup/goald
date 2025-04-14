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

const codeGenCLASSES codeGenLevel = 1
const codeGenSPECS codeGenLevel = 2
const codeGenUTILS codeGenLevel = 3
const dirtyFILENAME = "dirty"

// this function shows that our server, when run in dev with the right arguments,
// can be used as a development server, generating code for us
func (thisServer *server) runCodeGen(srcdir string, level codeGenLevel, webdir, nativedir string, regen bool, bindir string) {
	switch level {
	case codeGenCLASSES:
		start := time.Now()

		// TODO optimize with go routines here (?)

		// we're making all the databases globally accessible
		thisServer.generateDatabasesList(srcdir)

		// generating the classes and the packages that register them, and make the corresponding business objects "importable"
		codeChanged := thisServer.generateAllClasses(srcdir, ".", false, map[packageName]map[className]*classCore{}, regen)

		// saving the dirty state
		utils.WriteToFile(fmt.Sprintf("%t", codeChanged), bindir, dirtyFILENAME)

		slog.Info(fmt.Sprintf("done generating the DB & BO registries in %s", time.Since(start)))
		os.Exit(0)

	case codeGenSPECS:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build a static representation of our BOs
		codeChanged := thisServer.generateAllObjectSpecs(srcdir, regen)

		// saving the dirty state
		utils.WriteToFile(fmt.Sprintf("%t", codeChanged), bindir, dirtyFILENAME)

		slog.Info(fmt.Sprintf("done generating the BO specs in %s", time.Since(start)))
		os.Exit(0)

	case codeGenUTILS:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build utils for our BOs,
		// that should help us avoid using the `reflect` package at runtime;
		codeChanged := thisServer.generateAllObjectValueMappers(srcdir, ".", regen)

		// codegen in the webapp! and / or the native app
		thisServer.generateAllClientAppModels(webdir, regen, true)
		thisServer.generateAllClientAppModels(nativedir, regen, false)

		// saving the dirty state
		utils.WriteToFile(fmt.Sprintf("%t", codeChanged), bindir, dirtyFILENAME)

		slog.Info(fmt.Sprintf("done generating the BO utils & models in %s", time.Since(start)))
		os.Exit(0)

	default:
		utils.Panicf("Not handling to code generation level: %d", level)
	}
}

// ------------------------------------------------------------------------------------------------
// Utilities
// ------------------------------------------------------------------------------------------------
