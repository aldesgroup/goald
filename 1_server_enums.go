// ------------------------------------------------------------------------------------------------
// Here are the enums used for building business object classes
// ------------------------------------------------------------------------------------------------
package goald

// ------------------------------------------------------------------------------------------------
// the environment type for the currently running app
// ------------------------------------------------------------------------------------------------

// envType represents the type of environment we're running the app in
type envType int

const (
	envTypeDEV  envType = -1
	envTypeTEST envType = 1
	envTypePROD envType = 2
)

var envTypes = map[int]string{
	int(envTypeDEV):  "dev",
	int(envTypeTEST): "test",
	int(envTypePROD): "prod",
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

func envTypeFrom(value string) envType {
	for eT, label := range envTypes {
		if label == value {
			return envType(eT)
		}
	}

	return 0
}
