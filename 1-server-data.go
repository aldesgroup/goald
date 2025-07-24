// ------------------------------------------------------------------------------------------------
// The code here is about the loading of data for the server to work properly.
// Data loaders - which are functions with the same signature - must be registered through a
// function defined here.
// Beware: data loaders can be invoked in parallel, so there should not be any dependency between
// them. If you have some logic about several data loading bits, you must implement it within
// one data loader.
// ------------------------------------------------------------------------------------------------
package goald

import (
	"fmt"
	"log/slog"
	"sync"

	core "github.com/aldesgroup/corego"
)

// the data loader type
type dataLoader func(serverCtx BloContext, params map[string]string) error

// the main data loading function
func (thisServer *server) loadData(migrationPhase bool) {
	// handling the synchronization of the loaders
	wg := new(sync.WaitGroup)

	// gathering the errors (hopefully none)
	errors := map[string]error{}
	errorMx := new(sync.Mutex)

	// launching all the registred data loaders in parallel!
	loaders := core.IfThenElse(migrationPhase, dataLoaderRegistry.migrationLoaders, dataLoaderRegistry.appServerLoaders)
	for fnName, dataLoadingFn := range loaders {
		slog.Debug(fmt.Sprintf("Starting runner: %s", fnName))
		wg.Add(1)
		go func(fnNameArg string, dataLoadingFnArg dataLoader) {
			defer wg.Done()
			defer RecoverError("Error while running data loader '%s'", fnNameArg)

			if errLoad := dataLoadingFnArg(thisServer, thisServer.config.commonPart().DataLoaders[fnNameArg]); errLoad != nil {
				errorMx.Lock()
				errors[fnNameArg] = errLoad
				errorMx.Unlock()
			}
		}(fnName, dataLoadingFn)
	}

	// waiting for the last loader to finish
	wg.Wait()

	// for now, only logging the errors, but we might end up stopping the server in case of error
	for fnName, errLoad := range errors {
		slog.Error(fmt.Sprintf("Error in data loader '%s': %s", fnName, errLoad))
	}
}
