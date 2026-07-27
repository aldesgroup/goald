// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	_ "github.com/aldesgroup/goald/_include/iot/model"
	iot "github.com/aldesgroup/goald/features/iot/class"
)

func init() {
	g.In("goald").
		Register(iot.ClassForDevice("features/iot", "2026-07-22T13:06:37+02:00")).
		Register(iot.ClassForDeviceBootstrap("features/iot", "2026-07-22T13:06:37+02:00")).
		Register(iot.ClassForDeviceBootstrapPayload("features/iot", "2026-07-22T13:06:37+02:00")).
		Register(iot.ClassForDeviceBootstrapRequest("features/iot", "2026-07-22T13:06:37+02:00")).
		Register(iot.ClassForDeviceLinkRequest("features/iot", "2026-07-22T13:06:37+02:00"))
}
