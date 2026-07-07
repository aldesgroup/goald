// ------------------------------------------------------------------------------------------------
// Here is the global code generation routine
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
)

type codeGenLevel int
type packageName string

const codeGenCLASSES codeGenLevel = 1
const codeGenMODELS codeGenLevel = 2
const codeGenUTILS codeGenLevel = 3
const codeGenCHECK codeGenLevel = 4
const dirtyFILENAME = "dirty"

// this function shows that our server, when run in dev with the right arguments,
// can be used as a development server, generating code for us
func (thisServer *server) runCodeGen(cgp *codegenParams) {
	switch level := codeGenLevel(cgp.codegen); level {
	case codeGenCLASSES:
		start := time.Now()

		// TODO optimize with go routines here (?)

		// we're making all the databases globally accessible
		thisServer.generateDatabasesList(cgp.srcdir)

		// generating the classes and the packages that register them, and make the corresponding business objects "importable"
		codeChanged := thisServer.generateAllClasses(cgp.srcdir, ".", false, map[packageName]map[className]*classCore{}, cgp.regen)

		// saving the dirty state
		core.WriteToFile(fmt.Sprintf("%t", codeChanged), cgp.bindir, dirtyFILENAME)

		thisServer.Info(fmt.Sprintf("done generating the DB & BO registries in %s", time.Since(start)))

	case codeGenMODELS:
		start := time.Now()

		// now, using the `reflect` package, we can "easily" build a static representation of our BOs
		codeChanged := thisServer.generateAllObjectModels(cgp.srcdir, cgp.regen)

		// saving the dirty state
		core.WriteToFile(fmt.Sprintf("%t", codeChanged), cgp.bindir, dirtyFILENAME)

		thisServer.Info(fmt.Sprintf("done generating the BO models in %s", time.Since(start)))

	case codeGenUTILS:
		start := time.Now()

		// now, using the models, we can generate useful utils
		codeChanged := thisServer.generateAllObjectValueMappers(cgp.srcdir, ".", cgp.regen)

		// saving the dirty state
		core.WriteToFile(fmt.Sprintf("%t", codeChanged), cgp.bindir, dirtyFILENAME)

		// codegen in the webapp! and / or the native app
		thisServer.generateAllClientAppModels(cgp.webdir, cgp.regen, true)
		thisServer.generateAllClientAppModels(cgp.nativedir, cgp.regen, false)

		// generating the doc for the API
		if cgp.docpath != "" {
			srcdirs := append([]string{cgp.srcdir}, strings.Split(cgp.othersrcdirs, ",")...)
			thisServer.generateOpenAPIDoc(srcdirs, cgp.docpath, cgp.regen, cgp.servers)
		}

		thisServer.Info(fmt.Sprintf("done generating the BO utils, client models & API doc in %s", time.Since(start)))

	case codeGenCHECK:
		// at the end, we check the code is fine
		thisServer.runCodeChecks()

	default:
		core.PanicMsg("Not handling to code generation level: %d", level)
	}
}
