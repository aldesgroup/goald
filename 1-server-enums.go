// ------------------------------------------------------------------------------------------------
// Here are the enums used for building business object classes
// ------------------------------------------------------------------------------------------------
package goald

import (
	"strings"

	core "github.com/aldesgroup/corego"
)

// ------------------------------------------------------------------------------------------------
// the environment type for the currently running app
// ------------------------------------------------------------------------------------------------

// envType represents the type of environment we're running the app in
type envType int

const (
	envTypeLOCAL      envType = -2
	envTypeSANDBOX    envType = -1
	envTypeSTAGING    envType = 1
	envTypePRODUCTION envType = 2
)

var envTypes = map[int]string{
	int(envTypeLOCAL):      "LOCAL",
	int(envTypeSANDBOX):    "SANDBOX",
	int(envTypeSTAGING):    "STAGING",
	int(envTypePRODUCTION): "PRODUCTION",
}

func (thisEnvType envType) String() string {
	return envTypes[int(thisEnvType)]
}

// Val helps implement the IEnum interface
func (thisEnvType envType) Val() int {
	return int(thisEnvType)
}

// Values helps implement the IEnum interface
func (thisEnvType envType) Values() map[int]string {
	return envTypes
}

func envTypeValFrom(value string) envType {
	for eT, label := range envTypes {
		if label == value {
			return envType(eT)
		}
	}

	core.PanicMsg("There's no env type named '%s'. Possible values are: %s",
		value, strings.Join(core.GetSortedValues(envTypes), ", "))
	return 0
}
