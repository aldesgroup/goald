// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	iot "github.com/aldesgroup/goald/_include/iot/class"
)

func init() {
	g.In("goald").
		Register(iot.ClassForDevice("features/iot", "2026-07-27T16:36:01+02:00")).
		Register(iot.ClassForDeviceBootstrap("features/iot", "2026-07-27T16:36:01+02:00")).
		Register(iot.ClassForDeviceBootstrapPayload("features/iot", "2026-07-27T16:36:01+02:00")).
		Register(iot.ClassForDeviceBootstrapRequest("features/iot", "2026-07-27T16:36:01+02:00")).
		Register(iot.ClassForDeviceLinkRequest("features/iot", "2026-07-27T16:36:01+02:00"))
}
