package main

import (
	"strings"
)

func (o *OutputType) UpdateVlan(inputData InputData) {
	vlanList := []VlanType{}
	if strings.Contains(o.Switch.Type, TOR) {
		// TOR Switch with all matched vlans
		for _, supernet := range inputData.Supernets {
			vlanItem := VlanType{}
			if supernet.VlanID != 0 {
				vlanItem.VlanName = supernet.Name
				vlanItem.VlanID = supernet.VlanID
				vlanItem.GroupID = supernet.GroupID
				vlanItem.Mtu = JUMBOMTU
				if supernet.Shutdown {
					vlanItem.Shutdown = true
				}
				if supernet.IPv4.SwitchSVI {
					for _, ipv4 := range supernet.IPv4.Assignment {
						if ipv4.Name == VIPGATEWAY {
							vlanItem.VIPAddress = ipv4.IP
						} else if ipv4.Name == o.Switch.Type {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
						}
					}
				}
				vlanList = append(vlanList, vlanItem)
			}
		}
	} else if strings.Contains(o.Switch.Type, BMC) {
		// BMC Switch only have Unused and BMC Vlan (no VIP)
		for _, supernet := range inputData.Supernets {
			vlanItem := VlanType{}
			if supernet.GroupID == UNUSED || supernet.GroupID == BMC {
				vlanItem.VlanName = supernet.Name
				vlanItem.VlanID = supernet.VlanID
				vlanItem.GroupID = supernet.GroupID
				vlanItem.Mtu = DefaultMTU
				if supernet.Shutdown {
					vlanItem.Shutdown = true
				}
				if supernet.IPv4.SwitchSVI {
					for _, ipv4 := range supernet.IPv4.Assignment {
						if ipv4.Name == o.Switch.Type {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
						}
					}
				}
				vlanList = append(vlanList, vlanItem)
			}
		}
	}

	// Convert Vlan List to Map
	// vlanMap := map[string][]VlanType{}
	// for _, vlanObj := range vlanList {
	// 	groupId := vlanObj.GroupID
	// 	vlanMap[groupId] = append(vlanMap[groupId], vlanObj)
	// }
	o.Vlans = vlanList
}
