package server

import (
	"sync"

	"github.com/aldesgroup/goald"
	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/features/iot"
)

// ----------------------------------------------------------------------------
// CRUD part of the device management - simulated for now
// TODO remove this ultimately
// ----------------------------------------------------------------------------

var devices []*iot.Device
var devicesMx sync.Mutex

func init() {
	devices = []*iot.Device{}
}

func addDevice(model, serial string, user goald.IUser) *iot.Device {
	devicesMx.Lock()
	defer devicesMx.Unlock()

	device := &iot.Device{
		Status:          iot.DeviceStatusLINKED,
		Model:           model,
		Serial:          serial,
		AssociatedUsers: []goald.IUser{user},
	}

	// saving it
	// TODO obviously, we'll have to change this later
	devices = append(devices, device)

	return device
}

func getDevice(serial string) *iot.Device {
	devicesMx.Lock()
	defer devicesMx.Unlock()

	for _, device := range devices {
		if device.Serial == serial {
			return device
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Linking the current user to the current device, creating it on the fly
// if needed
// ----------------------------------------------------------------------------

func doLinkDevice(bloContext g.BloContext, model, serial string, userID int) (*iot.Device, error) {
	// does this user exist?
	user := iot.GetUser(userID)
	if user == nil {
		return nil, g.Error("No user found with id '%d'", userID)
	}

	// controlling the serial & the model
	// TODO we'll have to somehow plugin some stuff here, prolly coming from emerald
	if model != "SENSAIR" {
		return &iot.Device{Status: iot.DeviceStatusFORBIDDEN}, nil
	}
	if len(serial) != 13 {
		return &iot.Device{Status: iot.DeviceStatusBADxDEVICExSERIAL}, nil
	}

	// retrieving the device from the DB (mocked, for now)
	device := getDevice(serial)
	if device != nil {
		println(device.Status == iot.DeviceStatusLINKED)
		println(device.IsUserAssociated(user))
		if device.Status == iot.DeviceStatusLINKED && device.IsUserAssociated(user) {
			device.Status = iot.DeviceStatusALREADYxLINKED
			// TODO then save this
		}

		return device, nil
	}

	// creating a new device since we didn't have this one yet
	device = addDevice(model, serial, user)

	return device, nil
}
