// Generated file, do not edit!
package accessmgt

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/accessmgt/model"
	accessmgt "github.com/aldesgroup/goald/features/accessmgt/class"
)

func init() {
	g.In("goald").
		Register(accessmgt.ClassForUser("features/accessmgt", "2026-07-22T13:06:37+02:00")).
		Register(accessmgt.ClassForUserGroup("features/accessmgt", "2026-07-22T14:13:39+02:00"))
}
