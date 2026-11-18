// Generated file, do not edit!
package accessmgt

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/accessmgt/model"
	"github.com/aldesgroup/goald/_include/accessmgt/source"
)

func init() {
	g.In("goald").
		Register(source.ForBaseUserProfile("features/accessmgt", "2026-11-16T14:07:56+01:00")).
		Register(source.ForPermission("features/accessmgt", "2026-11-13T23:02:13+01:00")).
		Register(source.ForPermissionOwner("features/accessmgt", "2026-11-13T23:03:50+01:00")).
		Register(source.ForUser("features/accessmgt", "2026-11-18T13:53:20+01:00")).
		Register(source.ForUserGroup("features/accessmgt", "2026-11-18T09:01:34+01:00"))
}
