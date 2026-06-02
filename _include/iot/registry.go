// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/iot/model"
	iot "github.com/aldesgroup/goald/features/iot/class"
)

func init() {
	g.In("goald").
		Register(iot.ClassForDevice("features/iot", "2026-04-28T15:07:04+02:00")).
		Register(iot.ClassForDeviceBootstrap("features/iot", "2026-04-28T15:10:50+02:00")).
		Register(iot.ClassForDeviceBootstrapPayload("features/iot", "2026-04-28T15:10:03+02:00")).
		Register(iot.ClassForDeviceBootstrapRequest("features/iot", "2026-04-28T15:09:35+02:00")).
		Register(iot.ClassForDeviceLinkRequest("features/iot", "2026-04-28T23:50:17+02:00"))
}
