// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"os"

	core "github.com/aldesgroup/corego"
	"sigs.k8s.io/yaml"
)

// ------------------------------------------------------------------------------------------------
// Useful structs
// ------------------------------------------------------------------------------------------------

type IServerConfig interface {
	IBaseConfig
	CustomConfig() ICustomConfig // the applicative, custom part of the config
}

type IBaseConfig interface {
	base() *serverConfig // the common, generic part of the config
}

type ICustomConfig interface {
	// nothing for now
}

func NewBaseConfig() *serverConfig {
	return &serverConfig{}
}

type serverConfig struct {
	// HTTP        *httpConfig
	EnvType     string
	Databases   []*dbConfig
	DataLoaders map[string]map[string]string

	// technical props
	envTypeVal core.EnvType
}

type DatabaseID string

// type httpConfig struct {
// 	// Port int
// 	// ApiPath      string
// 	StaticRoutes []*staticRouteConfig
// }

type staticRouteConfig struct {
	For       string
	ServeFile string
	ServeDir  string
}

type dbConfig struct {
	DbID      DatabaseID
	DbType    databaseType
	DbName    string
	DbHost    string
	DbPort    int
	User      string
	Password  string
	MakeExist bool
}

// ------------------------------------------------------------------------------------------------
// Config reading
// ------------------------------------------------------------------------------------------------

var configObj IServerConfig

func RegisterConfig(cfgObj IServerConfig) {
	// doing this only once
	if configObj == nil {
		configObj = cfgObj
	}
}

func readAndCheckConfig(fromPath string) IServerConfig {
	// Do we have a configuration object ready?
	if configObj == nil {
		core.PanicMsg("No configuration object (implementing IServerConfig) has been registered!")
	}

	// Reading the config file into bytes
	yamlBytes, errRead := os.ReadFile(fromPath)
	core.PanicMsgIfErr(errRead, "Could not read config file at path '%s'", fromPath)

	// YAML -> JSON transformation, because JSON unmarshalling is better
	jsonBytes, errJson := yaml.YAMLToJSON(yamlBytes)
	core.PanicMsgIfErr(errJson, "Could not convert YAML to JSON '%s'", fromPath)
	// Unmarshalling the YAML file
	core.PanicMsgIfErr(errRead, "Could not read config file at path '%s'", fromPath)
	core.PanicMsgIfErr(json.Unmarshal(jsonBytes, configObj),
		"Could not unmarshal the config file at path '%s'", fromPath)

	// controlling the common config
	config := configObj.base()

	// Parsing the env type
	config.envTypeVal = core.EnvTypeValFrom(config.EnvType)

	return configObj
}

func (thisConf *serverConfig) base() *serverConfig {
	return thisConf
}
