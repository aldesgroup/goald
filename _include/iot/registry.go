// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/iot/model"
	"github.com/aldesgroup/goald/_include/iot/source"
)

func init() {
	g.In("goald").
		Register(source.ForDevice("features/iot", "2026-08-03T16:50:33+02:00")).
		Register(source.ForDeviceBootstrap("features/iot", "2026-08-03T16:50:33+02:00")).
		Register(source.ForDeviceBootstrapPayload("features/iot", "2026-08-03T16:50:33+02:00")).
		Register(source.ForDeviceBootstrapRequest("features/iot", "2026-08-03T16:50:33+02:00")).
		Register(source.ForDeviceLinkRequest("features/iot", "2026-08-03T16:50:33+02:00"))
}
