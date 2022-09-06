package main

import (
	"sort"
	"strings"
)

func (o *OutputType) parseVlanObj(i *InterfaceFrameworkType) {
	for _, vlanItem := range i.Vlan {
		for _, segment := range *o.Supernets {
			if segment.Group == vlanItem.Group {
				vlanItem.VlanName = segment.Name
				vlanItem.Shutdown = segment.Shutdown
				vlanItem.VlanID = segment.VlanID
				IPAssignmentName := replaceTORXName(vlanItem.IPAddress, o.Device.Type)
				for _, ipAssign := range segment.IPAssignment {
					// Vlan Name Match
					if strings.Contains(ipAssign.Name, IPAssignmentName) {
						vlanItem.IPAddress = ipAssign.IPAddress
					}
					// VIP Name Match
					if strings.Contains(ipAssign.Name, vlanItem.Vip.VIPAddress) {
						vlanItem.Vip.VIPAddress = ipAssign.IPAddress
					}
				}
				o.Vlan = append(o.Vlan, vlanItem)
			}
		}
	}
	// Increasing Sort by VlanID
	sort.Slice(o.Vlan, func(i, j int) bool {
		return o.Vlan[i].VlanID < o.Vlan[j].VlanID
	})
}
