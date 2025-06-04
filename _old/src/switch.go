package main

import "strings"

func (i *InputData) createDeviceTypeMap() map[string][]SwitchType {
	DeviceTypeMap := map[string][]SwitchType{}
	for _, switchItem := range i.Switches {
		typeUpper := strings.ToUpper(switchItem.Type)
		if strings.Contains(typeUpper, BORDER) {
			DeviceTypeMap[BORDER] = append(DeviceTypeMap[BORDER], switchItem)
		} else if strings.Contains(typeUpper, TOR) {
			DeviceTypeMap[TOR] = append(DeviceTypeMap[TOR], switchItem)
		} else if strings.Contains(typeUpper, BMC) {
			DeviceTypeMap[BMC] = append(DeviceTypeMap[BMC], switchItem)
		} else if strings.Contains(typeUpper, MUX) {
			DeviceTypeMap[MUX] = append(DeviceTypeMap[MUX], switchItem)
		}
	}
	return DeviceTypeMap
}

func (o *OutputType) UpdateSwitch(switchItem SwitchType, switchType string, DeviceTypeMap map[string][]SwitchType) {
	o.SwitchUplink = DeviceTypeMap[BORDER]

	if switchType == TOR {
		o.Switch = switchItem
		o.SwitchDownlink = DeviceTypeMap[MUX]
		o.SwitchBMC = DeviceTypeMap[BMC]
		for _, torItem := range DeviceTypeMap[TOR] {
			if switchItem != torItem {
				o.SwitchPeer = append(o.SwitchPeer, torItem)
			}
		}
	} else if switchType == BMC {
		o.Switch = switchItem
	}
}
