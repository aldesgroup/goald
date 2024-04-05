// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"os"

	"github.com/aldesgroup/goald/features/utils"
	"sigs.k8s.io/yaml"
)

// ------------------------------------------------------------------------------------------------
// Useful structs
// ------------------------------------------------------------------------------------------------

type IServerConfig interface {
	ICommonConfig
	CustomConfig() ICustomConfig // the applicative, custom part of the config
}

type ICommonConfig interface {
	commonPart() *serverConfig // the common, generic part of the config
}

type ICustomConfig interface {
	// nothing for now
}

func NewCommonConfig() *serverConfig {
	return &serverConfig{}
}

type serverConfig struct {
	Env       string
	HTTP      *httpConfig
	Databases []*dbConfig

	// technical props
	envAsType envType
}

type DatabaseID string

type httpConfig struct {
	Port         int
	ApiPath      string
	StaticRoutes []*staticRouteConfig
}

type staticRouteConfig struct {
	For       string
	ServeFile string
	ServeDir  string
}

type dbConfig struct {
	DbID     DatabaseID
	DbType   databaseType
	DbName   string
	DbPort   string
	User     string
	Password string
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
		utils.Panicf("No configuration object (implementing IServerConfig) has been registered!")
	}

	// Reading the config file into bytes
	yamlBytes, errRead := os.ReadFile(fromPath)
	utils.PanicErrf(errRead, "Could not read config file at path '%s'", fromPath)

	// YAML -> JSON transformation, because JSON unmarshalling is better
	jsonBytes, errJson := yaml.YAMLToJSON(yamlBytes)
	utils.PanicErrf(errJson, "Could not convert YAML to JSON '%s'", fromPath)

	// Unmarshalling the YAML file
	utils.PanicErrf(errRead, "Could not read config file at path '%s'", fromPath)
	utils.PanicErrf(json.Unmarshal(jsonBytes, configObj),
		"Could not unmarshal the config file at path '%s'", fromPath)

	// controlling the common config
	config := configObj.commonPart()

	// Checking the env type
	if config.envAsType = envTypeFrom(config.Env); config.envAsType == 0 {
		utils.Panicf("the 'Env' config item (\"%s\") is not set, or not one of these values: dev, test, prod",
			config.Env)
	}

	return configObj
}

func (thisConf *serverConfig) commonPart() *serverConfig {
	return thisConf
}
