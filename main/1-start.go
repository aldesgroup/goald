package main

import (
	"github.com/aldesgroup/goald"

	// sourcing other features
	_ "github.com/aldesgroup/goald/features/i18n"
)

func main() {
	goald.NewServer().Start()
}
