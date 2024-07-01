// Generated file, do not edit!
package main

import (
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/i18n"
)

func init() {
	g.In("goald").
		Register(func() any { return &i18n.TestObject{} }, "features/i18n", "2024-06-27T00:16:16+02:00", func() any { return []*i18n.TestObject{} }).
		Register(func() any { return &i18n.Translation{} }, "features/i18n", "2024-06-17T09:12:08+02:00", func() any { return []*i18n.Translation{} }).
		Register(func() any { return &i18n.TranslationUrlParams{} }, "features/i18n", "2024-06-26T13:11:52+02:00", func() any { return []*i18n.TranslationUrlParams{} })
}
