// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/iot/model"
	"github.com/aldesgroup/goald/_include/iot/source"
)

func init() {
	g.In("goald").
		Register(source.ForDevice("features/iot", "2026-09-28T15:49:28+02:00")).
		Register(source.ForDeviceBootstrap("features/iot", "2026-09-28T15:35:35+02:00")).
		Register(source.ForDeviceBootstrapPayload("features/iot", "2026-09-28T15:35:35+02:00")).
		Register(source.ForDeviceBootstrapRequest("features/iot", "2026-09-28T15:35:35+02:00")).
		Register(source.ForDeviceLinkRequest("features/iot", "2026-09-28T15:50:15+02:00"))
}
