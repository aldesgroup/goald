// Generated file, do not edit!
package accessmgt

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/accessmgt/model"
	"github.com/aldesgroup/goald/_include/accessmgt/source"
)

func init() {
	g.In("goald").
		Register(source.ForUser("features/accessmgt", "2026-07-30T11:35:48+02:00")).
		Register(source.ForUserGroup("features/accessmgt", "2026-07-30T11:35:57+02:00"))
}
