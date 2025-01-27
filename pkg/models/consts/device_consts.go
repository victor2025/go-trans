package consts

/**

@author: victor2022
@since: 2025/1/27
*/

type DeviceMode string

const (
	DeviceScanIdentityParam            = "go-trans"
	SendMode                DeviceMode = "send"
	ReceiveMode             DeviceMode = "receive"
	SelfMode                DeviceMode = "self"
)
