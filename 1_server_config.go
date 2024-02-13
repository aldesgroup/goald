// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"os"

	"github.com/tidwall/jsonc"
)

// ------------------------------------------------------------------------------------------------
// Useful structs
// ------------------------------------------------------------------------------------------------

type IServerConfig interface {
	ICommonConfig
	CustomPart() ICustomConfig // the applicative, custom part of the config
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
	fileBytes, errRead := os.ReadFile(fromPath)
	panicErrf(errRead, "Could not read config file at path '%s'", fromPath)

	// Unmarshalling the JSONC file
	panicErrf(errRead, "Could not read config file at path '%s'", fromPath)
	panicErrf(json.Unmarshal(jsonc.ToJSON(fileBytes), intoConfigObj),
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
