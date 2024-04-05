package main

import g "github.com/aldesgroup/goald"

func init() {
	g.RegisterConfig(&goaldConfig{
		ICommonConfig: g.NewCommonConfig(),
	})
}

type goaldConfig struct {
	g.ICommonConfig `json:"Common"`
	// no custom config for now
}

func (thisConf *goaldConfig) CustomConfig() g.ICustomConfig {
	return nil // no custom config for now
}
