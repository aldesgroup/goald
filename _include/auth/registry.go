// Generated file, do not edit!
package auth

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/auth/model"
	auth "github.com/aldesgroup/goald/features/auth/class"
)

func init() {
	g.In("goald").
		Register(auth.ClassForUser("features/auth", "2026-06-30T08:26:30+02:00"))
}
