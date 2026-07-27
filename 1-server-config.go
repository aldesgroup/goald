// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/logging"
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
	AppName string
	AppDesc string
	Port    int
	EnvType string // LOCAL, SANDBOX, STAGING, PRODUCTION
	Logging *struct {
		Level  string              // debug, info, warn, error
		Type   logging.LoggingType // text or json
		FNames bool                // whether to include function names in the logs

		// tech props
		resolvedLogLevel slog.Level
	}
	DBServers   map[string]*dbconn.DbServerConfig
	DataLoaders map[string]map[string]string
	Version     string

	// technical props
	resolvedEnvType core.EnvType
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
	core.PanicMsgIfErr(json.Unmarshal(jsonBytes, configObj),
		"Could not unmarshal the config file at path '%s'", fromPath)

	// controlling the common config
	config := configObj.base()

	// Parsing the env type
	config.resolvedEnvType = core.EnvTypeValFrom(config.EnvType)

	// Parsing the log level
	core.PanicMsgIf(config.Logging == nil, "No \"Logging\" section configured!")
	config.Logging.resolvedLogLevel = slog.LevelInfo
	switch strings.ToUpper(config.Logging.Level) {
	case "TRACE":
		config.Logging.resolvedLogLevel = logging.LevelTrace
	case slog.LevelDebug.String():
		config.Logging.resolvedLogLevel = slog.LevelDebug
	case slog.LevelInfo.String():
		config.Logging.resolvedLogLevel = slog.LevelInfo
	case slog.LevelWarn.String():
		config.Logging.resolvedLogLevel = slog.LevelWarn
	case slog.LevelError.String():
		config.Logging.resolvedLogLevel = slog.LevelError
	default:
		core.PanicMsg("Invalid log level '%s' in config file", config.Logging.Level)
	}

	// Checking the logging type
	if !core.InSlice([]logging.LoggingType{logging.LoggingTypeTEXT, logging.LoggingTypeJSON}, config.Logging.Type) {
		core.PanicMsg("Invalid logging type '%s' in config file", config.Logging.Type)
	}

	// controlling the DB servers
	for dbID, dbServer := range config.DBServers {
		if dbServer.Type == "" {
			core.PanicMsg("DB server '%s' has no type defined", dbID)
		}
		if !core.InSlice(allDbTypes, dbServer.Type) {
			core.PanicMsg("DB server '%s' has an invalid type '%s'", dbID, dbServer.Type)
		}
		if dbServer.Host == "" {
			core.PanicMsg("DB server '%s' has no host defined", dbID)
		}
		if dbServer.Port <= 0 {
			core.PanicMsg("DB server '%s' has no port defined", dbID)
		}
		if dbServer.Database == "" {
			core.PanicMsg("DB server '%s' has no database name defined", dbID)
		}
		for schemaName, schemaConfig := range dbServer.Schemas {
			schemaConfig.Name = schemaName
			schemaConfig.DbServer = dbServer
		}
	}

	return configObj
}

func (thisConf *serverConfig) base() *serverConfig {
	return thisConf
}
