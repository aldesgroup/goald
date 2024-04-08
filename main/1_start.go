package main

import (
	"github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/features/i18n"
	// _ "github.com/aldesgroup/goald/features/other"
)

func main() {
	goald.NewServer().Start()
}
