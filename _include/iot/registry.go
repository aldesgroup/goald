// Generated file, do not edit!
package iot

import (
	g "github.com/aldesgroup/goald"
	iot "github.com/aldesgroup/goald/features/iot/class"
)

func init() {
	g.In("goald").
		Register(iot.ClassForDevice("features/iot", "2026-03-10T11:51:26+01:00")).
		Register(iot.ClassForDeviceBootstrap("features/iot", "2026-03-11T15:57:41+01:00")).
		Register(iot.ClassForDeviceBootstrapPayload("features/iot", "2026-03-12T14:25:58+01:00")).
		Register(iot.ClassForDeviceBootstrapRequest("features/iot", "2026-03-05T16:19:22+01:00")).
		Register(iot.ClassForDeviceLinkRequest("features/iot", "2026-03-10T11:25:38+01:00"))
}
