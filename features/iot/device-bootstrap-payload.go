package iot

type DeviceBootstrapPayload struct {
	Serial    string
	Model     string
	Timestamp int64
	Nonce     string
}
