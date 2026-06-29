package main

import (
	"github.com/aldesgroup/goald"

	// we're having a "server" here only for code-generation in the following packages:
	_ "github.com/aldesgroup/goald/_include/auth"
	_ "github.com/aldesgroup/goald/_include/i18n"
	_ "github.com/aldesgroup/goald/_include/iot"
)

func main() {
	goald.NewServer().Start()
}
