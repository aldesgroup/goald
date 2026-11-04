// ------------------------------------------------------------------------------------------------
// Here we configure the the server
// ------------------------------------------------------------------------------------------------
package goald

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"strings"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/auth"
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
	DBServers       map[string]*dbconn.DbServerConfig
	Auth            map[string]*auth.ProviderConfig // authentication realms, e.g. "internal" / "external" -> their provider config; entirely optional
	AuthTokenHeader string                          // header carrying the user's bearer token when a gateway keeps "Authorization" for its own token; defaults to "Authorization"
	DataLoaders     map[string]map[string]string
	Version         string

	// technical props
	resolvedEnvType core.EnvType
}

// ------------------------------------------------------------------------------------------------
// Config reading
// ------------------------------------------------------------------------------------------------

var configObj IServerConfig

func RegisterConfig(cfgObj IServerConfig, resolver ISecretResolver) {
	// doing this only once
	if configObj == nil {
		configObj = cfgObj
		secretResolver = resolver
	}
}

func readAndCheckConfig(fromPath string, resolveSecretValues bool) IServerConfig {
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

	// resolving the "secret:XXX" references before unmarshalling
	if resolveSecretValues {
		var tree any
		core.PanicMsgIfErr(json.Unmarshal(jsonBytes, &tree), "Could not parse the config file at path '%s'", fromPath)
		jsonBytes, errJson = json.Marshal(resolveSecrets(tree, path.Base(fromPath) == "conf-local.yaml"))
		core.PanicMsgIfErr(errJson, "Could not re-encode the config file at path '%s'", fromPath)
	}

	// Unmarshalling the YAML file
	core.PanicMsgIfErr(json.Unmarshal(jsonBytes, configObj), "Could not unmarshal the config file at path '%s'", fromPath)

	// controlling the base config
	baseConfig := configObj.base()

	// Parsing the env type
	baseConfig.resolvedEnvType = core.EnvTypeValFrom(baseConfig.EnvType)

	// Parsing the log level
	core.PanicMsgIf(baseConfig.Logging == nil, "No \"Logging\" section configured!")
	baseConfig.Logging.resolvedLogLevel = slog.LevelInfo
	switch strings.ToUpper(baseConfig.Logging.Level) {
	case "TRACE":
		baseConfig.Logging.resolvedLogLevel = logging.LevelTrace
	case slog.LevelDebug.String():
		baseConfig.Logging.resolvedLogLevel = slog.LevelDebug
	case slog.LevelInfo.String():
		baseConfig.Logging.resolvedLogLevel = slog.LevelInfo
	case slog.LevelWarn.String():
		baseConfig.Logging.resolvedLogLevel = slog.LevelWarn
	case slog.LevelError.String():
		baseConfig.Logging.resolvedLogLevel = slog.LevelError
	default:
		core.PanicMsg("Invalid log level '%s' in config file", baseConfig.Logging.Level)
	}

	// Checking the logging type
	if !core.InSlice([]logging.LoggingType{logging.LoggingTypeTEXT, logging.LoggingTypeJSON}, baseConfig.Logging.Type) {
		core.PanicMsg("Invalid logging type '%s' in config file", baseConfig.Logging.Type)
	}

	// controlling the DB servers
	for dbID, dbServer := range baseConfig.DBServers {
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

	// controlling the auth realms, if any - entirely optional; provider-specific requirements
	// (e.g. azuread's TenantID / Audience) are checked when the provider is used
	applyAuthEnvOverrides(baseConfig)
	for realmID, providerCfg := range baseConfig.Auth {
		if providerCfg.Type == "" {
			core.PanicMsg("Auth realm '%s' has no provider type defined", realmID)
		}
	}

	return configObj
}

func (thisConf *serverConfig) base() *serverConfig {
	return thisConf
}

// ------------------------------------------------------------------------------------------------
// Secret resolver definition & main secret resolution function
// ------------------------------------------------------------------------------------------------

const secretPrefix = "secret:"

var secretResolver ISecretResolver

type ISecretResolver interface {
	ResolveSecret(name string) string
}

// resolveSecrets replaces, in place, every string value starting with "secret:" by its resolved value
func resolveSecrets(node any, isLocal bool) any {
	switch typed := node.(type) {
	case map[string]any:
		for key, value := range typed {
			typed[key] = resolveSecrets(value, isLocal)
		}
	case []any:
		for i, value := range typed {
			typed[i] = resolveSecrets(value, isLocal)
		}
	case string:
		if name, found := strings.CutPrefix(typed, secretPrefix); found {
			// in local mode, we read the secret from the environment variable
			if isLocal {
				resolved := os.Getenv(name)
				if resolved == "" {
					core.PanicMsg("The value for secret '%s' in the .env.local file is empty!", name)
				}
				return resolved
			}

			// otherwise, we resolve it using the secret resolver
			return secretResolver.ResolveSecret(name)
		}
	}

	return node
}

// ------------------------------------------------------------------------------------------------
// Auth handling - for now
// ------------------------------------------------------------------------------------------------

// applyAuthEnvOverrides lets a deployment inject auth settings through env vars:
// AUTH_TOKEN_HEADER, and AUTH_<REALM>_<SETTING> for each realm declared in the config file.
func applyAuthEnvOverrides(config *serverConfig) {
	if value := os.Getenv("AUTH_TOKEN_HEADER"); value != "" {
		config.AuthTokenHeader = value
	}

	for realmID, providerCfg := range config.Auth {
		prefix := "AUTH_" + strings.ToUpper(realmID) + "_"
		settings := map[string]*string{
			"TENANT_ID": &providerCfg.TenantID,
			"AUDIENCE":  &providerCfg.Audience,
			"AUTHORITY": &providerCfg.Authority,
			"ISSUER":    &providerCfg.Issuer,
		}

		for suffix, target := range settings {
			if value := os.Getenv(prefix + suffix); value != "" {
				*target = value
			}
		}
	}
}
