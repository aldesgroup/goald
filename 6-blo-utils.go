package goald

import (
	"fmt"
	"slices"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/utils"
)

// ------------------------------------------------------------------------------------------------
// Handling slices of business objects
// ------------------------------------------------------------------------------------------------

// AddAndGetIDs adds all the given business objects to the cache and returns a map of all their IDs grouped by model name.
func (thisCache *BObjCache) AddAndGetIDs[BOTYPE IBusinessObject](bObjs []BOTYPE) (allBObjIDs map[utils.ModelName][]any) {
	allBObjIDs = make(map[utils.ModelName][]any)

	// iterating to build the result
	for _, bObj := range bObjs {
		// caching the object
		thisCache.Set(bObj)

		// gathering its ID along with the IDs of the objects of the same model
		allBObjIDs[bObj.GetModelName()] = append(allBObjIDs[bObj.GetModelName()], bObj.GetID())
	}

	return
}

// Makes a new slices with clones of the given business objects, and retrieves how they've being loaded
func cloneBOsWithLoading[BOTYPE IBusinessObject](model IBusinessObjectModel, bObjs []BOTYPE) (clones []BOTYPE, loadingConf ILoadingConfig, err error) {
	// nothing to do here
	if len(bObjs) == 0 {
		return nil, nil, nil
	}

	// getting one BO as a reference to retrieve how it has been loaded, so we can load the clones the same way
	ref := bObjs[0]
	withRelationships := []ILoadingConfig{}
	for _, loaded := range ref.getLoaded() {
		withRelationships = append(withRelationships, With(model.getRelationship(string(loaded))))
	}
	loadingConf = Load(model, withRelationships...)

	// checking all the given BOs have been loaded the same way, and clones for each of them
	clones = make([]BOTYPE, len(bObjs))
	for i, bObj := range bObjs {
		// not cloning a BO without an ID
		if bObj.GetID() <= 0 {
			return nil, nil, fmt.Errorf("cannot clone a business object without an ID (the one at index %d)", i)
		}

		// retrieving how the original BO has been loaded, so we can load the clone the same way
		if i > 0 && !bObj.getLoaded().equals(ref.getLoaded()) {
			return nil, nil, fmt.Errorf("the given business objects are not loaded the same way: %+v (0) vs %+v (%d)", ref.getLoaded(), bObj.getLoaded(), i)
		}

		// adding an empty BO of the same model with the same ID
		clones[i] = bObj.Clone(false, false).(BOTYPE)
	}

	return
}

// ------------------------------------------------------------------------------------------------
// Loading configs
// ------------------------------------------------------------------------------------------------

type ILoadingConfig interface {
	withParent(parent ILoadingConfig) ILoadingConfig
	getRelationshipsToLoad() []ILoadingConfig
	getLoadedRelationships() loadedRelationships
	getCurrentRelationship() *Relationship
	ToString(indent ...int) string
}

type loadingConfig struct {
	parent              ILoadingConfig
	relationshipsOwner  IBusinessObjectModel
	currentRelationship *Relationship
	relationshipsToLoad []ILoadingConfig
	relationshipsNames  loadedRelationships
}

func (cfg *loadingConfig) withParent(parent ILoadingConfig) ILoadingConfig {
	cfg.parent = parent
	return cfg
}

func (cfg *loadingConfig) getRelationshipsToLoad() []ILoadingConfig {
	return cfg.relationshipsToLoad
}

func (cfg *loadingConfig) getLoadedRelationships() loadedRelationships {
	return cfg.relationshipsNames
}

func (cfg *loadingConfig) getCurrentRelationship() *Relationship {
	return cfg.currentRelationship
}

func (cfg *loadingConfig) ToString(indent ...int) string {
	ind := 0
	if len(indent) > 0 {
		ind = indent[0]
	}
	prefix := strings.Repeat("  ", ind)

	result := ""

	if cfg.parent == nil {
		result = prefix + fmt.Sprintf("loading '%s' with:", cfg.relationshipsOwner.GetName())
	} else if ind == 0 {
		result = prefix + fmt.Sprintf("loading '%s' (from '%s#%s')%s",
			strings.Join(core.ToStrings(cfg.currentRelationship.getTargetModelNames()), ", "),
			cfg.relationshipsOwner.GetName(),
			cfg.currentRelationship.name,
			core.IfThenElse(len(cfg.relationshipsToLoad) > 0, " with:", "."),
		)
	}

	for _, subConfig := range cfg.relationshipsToLoad {
		result += "\n" + prefix + fmt.Sprintf("  - %s (%v)", subConfig.getCurrentRelationship().name, subConfig.getCurrentRelationship().getTargetModelNames())
		result += subConfig.ToString(ind + 1)
	}

	return result
}

func Load(model IBusinessObjectModel, with ...ILoadingConfig) ILoadingConfig {
	thisConfig := &loadingConfig{
		relationshipsOwner: model,
	}
	for _, subConfig := range with {
		thisConfig.relationshipsToLoad = append(thisConfig.relationshipsToLoad, subConfig.withParent(thisConfig))
		thisConfig.relationshipsNames = append(thisConfig.relationshipsNames, loadedRelationship(subConfig.getCurrentRelationship().GetName()))
	}
	slices.Sort(thisConfig.relationshipsNames)
	return thisConfig
}

func With(relationship *Relationship, with ...ILoadingConfig) ILoadingConfig {
	thisConfig := &loadingConfig{
		relationshipsOwner:  relationship.owner,
		currentRelationship: relationship,
	}
	for _, subConfig := range with {
		thisConfig.relationshipsToLoad = append(thisConfig.relationshipsToLoad, subConfig.withParent(thisConfig))
		thisConfig.relationshipsNames = append(thisConfig.relationshipsNames, loadedRelationship(subConfig.getCurrentRelationship().GetName()))
	}
	slices.Sort(thisConfig.relationshipsNames)
	return thisConfig
}

// ReadNoRelationship returns a loading config that loads only the direct relationships of the given business object type
func (thisModel *businessObjectModel) ReadNoRelationship() ILoadingConfig {
	return Load(thisModel)
}

// ReadWithFirstLayer returns a loading config that loads only the direct relationships of the given business object type
func (thisModel *businessObjectModel) ReadWithFirstLayer() ILoadingConfig {
	firstLayer := []ILoadingConfig{}
	for _, relationship := range core.GetSortedValues(thisModel.getRelationships()) {
		firstLayer = append(firstLayer, With(relationship))
	}
	return Load(thisModel, firstLayer...)
}

// ------------------------------------------------------------------------------------------------
// Keeping track of the loaded relationships for a business object
// ------------------------------------------------------------------------------------------------

type loadedRelationship string
type loadedRelationships []loadedRelationship

func (loaded loadedRelationships) equals(other loadedRelationships) bool {
	if len(loaded) != len(other) {
		return false
	}
	for i, rel := range loaded {
		if rel != other[i] {
			return false
		}
	}
	return true
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

func logObjs(bloCtx BloContext, bObjs []IBusinessObject, reasonForLogging string, args ...any) {
	if bloCtx.IsTraceEnabled() {
		objStrings := fmt.Sprintf("[%s]", strings.Join(core.MapFn(bObjs, func(bObj IBusinessObject) string {
			return fmt.Sprintf("%s (%p)", KeyFor(bObj), bObj)
		}), ", "))
		bloCtx.Trace("----------------------------------------------")
		bloCtx.Trace(fmt.Sprintf(reasonForLogging+": "+objStrings, args...))
		bloCtx.Trace("----------------------------------------------")
	}
}

func logIDs(bloCtx BloContext, bObjIDs []any, reasonForLogging string, args ...any) {
	if bloCtx.IsTraceEnabled() {
		objStrings := fmt.Sprintf("[%s]", strings.Join(core.MapFn(bObjIDs, func(id any) string {
			return fmt.Sprintf("%d", id)
		}), ", "))
		bloCtx.Trace("----------------------------------------------")
		bloCtx.Trace(fmt.Sprintf(reasonForLogging+": "+objStrings, args...))
		bloCtx.Trace("----------------------------------------------")
	}
}
