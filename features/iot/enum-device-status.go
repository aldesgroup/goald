package iot

type DeviceStatus int

const (
	DeviceStatusFORBIDDEN         DeviceStatus = -2 // The user can't be added to this device and / or this device can't be enrolled
	DeviceStatusBADxDEVICExSERIAL DeviceStatus = -1 // The serial number cannot help us identify a device type
	DeviceStatusLINKED            DeviceStatus = 1  // The user has successfully been associated with this device
	DeviceStatusALREADYxLINKED    DeviceStatus = 2  // The user was already associated with the given device (idempotence).
)

var deviceStatuses = map[int]string{
	int(DeviceStatusFORBIDDEN):         "forbidden",
	int(DeviceStatusBADxDEVICExSERIAL): "bad device serial",
	int(DeviceStatusLINKED):            "linked",
	int(DeviceStatusALREADYxLINKED):    "already linked",
}

func (thisDeviceStatus DeviceStatus) String() string {
	return deviceStatuses[int(thisDeviceStatus)]
}

// Val helps implement the IEnum interface
func (thisDeviceStatus DeviceStatus) Val() int {
	return int(thisDeviceStatus)
}

// Values helps implement the IEnum interface
func (thisDeviceStatus DeviceStatus) Values() map[int]string {
	return deviceStatuses
}
