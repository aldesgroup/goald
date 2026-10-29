// Generated file, do not edit!
package accessmgt

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/accessmgt/model"
	"github.com/aldesgroup/goald/_include/accessmgt/source"
)

func init() {
	g.In("goald").
		Register(source.ForAuthToken("features/accessmgt", "2026-10-05T14:24:33+02:00")).
		Register(source.ForLoginCredentials("features/accessmgt", "2026-09-30T16:54:17+02:00")).
		Register(source.ForUser("features/accessmgt", "2026-10-05T14:24:33+02:00")).
		Register(source.ForUserGroup("features/accessmgt", "2026-09-30T16:40:09+02:00"))
}
