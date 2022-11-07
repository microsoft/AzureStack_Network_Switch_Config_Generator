package main

import "strings"

func (i *InputData) createDeviceTypeMap() map[string][]SwitchType {
	DeviceTypeMap := map[string][]SwitchType{}
	for _, switchItem := range i.Switches {
		typeUpper := strings.ToUpper(switchItem.Type)
		if strings.Contains(typeUpper, BORDER) {
			DeviceTypeMap[UPLINK] = append(DeviceTypeMap[UPLINK], switchItem)
		} else if strings.Contains(typeUpper, TOR) {
			DeviceTypeMap[TOR] = append(DeviceTypeMap[TOR], switchItem)
		} else if strings.Contains(typeUpper, BMC) {
			DeviceTypeMap[BMC] = append(DeviceTypeMap[BMC], switchItem)
		} else if strings.Contains(typeUpper, MUX) {
			DeviceTypeMap[DOWNLINK] = append(DeviceTypeMap[DOWNLINK], switchItem)
		}
	}
	return DeviceTypeMap
}
