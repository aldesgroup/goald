// Generated file, do not edit!
package accessmgt

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/accessmgt/model"
	"github.com/aldesgroup/goald/_include/accessmgt/source"
)

func init() {
	g.In("goald").
		Register(source.ForUser("features/accessmgt", "2026-08-03T16:37:24+02:00")).
		Register(source.ForUserGroup("features/accessmgt", "2026-08-03T16:37:24+02:00"))
}
