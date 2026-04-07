// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/iot/model"
	iot "github.com/aldesgroup/goald/features/iot/class"
)

func init() {
	g.In("goald").
		Register(iot.ClassForDevice("features/iot", "2026-04-06T00:47:55+02:00")).
		Register(iot.ClassForDeviceBootstrap("features/iot", "2026-04-06T00:48:00+02:00")).
		Register(iot.ClassForDeviceBootstrapPayload("features/iot", "2026-04-06T00:48:05+02:00")).
		Register(iot.ClassForDeviceBootstrapRequest("features/iot", "2026-04-06T00:48:10+02:00")).
		Register(iot.ClassForDeviceLinkRequest("features/iot", "2026-04-06T00:48:14+02:00"))
}
