// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"os"

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
// Useful other types & constants
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Base implems
// ------------------------------------------------------------------------------------------------

func readAndCheckConfig(fromPath string, intoConfigObj IServerConfig) {
	// Reading the config file into bytes
	yamlBytes, errRead := os.ReadFile(fromPath)
	panicErrf(errRead, "Could not read config file at path '%s'", fromPath)

	// YAML -> JSON transformation, because JSON unmarshalling is better
	jsonBytes, errJson := yaml.YAMLToJSON(yamlBytes)
	panicErrf(errJson, "Could not convert YAML to JSON '%s'", fromPath)

	println(string(jsonBytes))

	// Unmarshalling the YAML file
	panicErrf(errRead, "Could not read config file at path '%s'", fromPath)
	panicErrf(json.Unmarshal(jsonBytes, intoConfigObj),
		"Could not unmarshal the config file at path '%s'", fromPath)

	// controlling the common config
	config := intoConfigObj.commonPart()

	// Checking the env type
	if config.envAsType = envTypeFrom(config.Env); config.envAsType == 0 {
		panicf("the 'Env' config item (\"%s\") is not set, or not one of these values: dev, test, prod",
			config.Env)
	}
}

func (thisConf *serverConfig) commonPart() *serverConfig {
	return thisConf
}
