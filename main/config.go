package main

import (
	g "github.com/aldesgroup/goald"
)

func init() {
	g.RegisterConfig(&goaldConfig{
		IBaseConfig: g.NewBaseConfig(),
	}, nil)
}

type goaldConfig struct {
	g.IBaseConfig `json:"base"`
	// no custom config for now
}

func (thisConf *goaldConfig) CustomConfig() g.ICustomConfig {
	return nil // no custom config for now
}
