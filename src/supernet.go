package main

import (
	"sort"
	"strings"
)

func (o *OutputType) UpdateVlanAndL3Intf(inputData InputData) {
	vlanList := []VlanType{}
	l3IntfMap := map[string]L3IntfType{}
	if strings.Contains(o.Switch.Type, TOR) {
		// TOR Switch with all matched vlans
		for idx, supernet := range inputData.Supernets {
			vlanItem := VlanType{}
			l3IntfItem := L3IntfType{}
			// Update Vlan ID but Switchless Deployment Skip Storage Vlan
			if supernet.IPv4.VlanID != 0 {
				if strings.Contains(strings.ToUpper(supernet.GroupName), strings.ToUpper(BMC)) {
					// BMC Vlan
					BMC_VlanID = supernet.IPv4.VlanID
				} else if strings.EqualFold(supernet.GroupName, Compute_NativeVlanName) {
					// Management/Infra Vlan
					Compute_NativeVlanID = supernet.IPv4.VlanID
				} else if strings.EqualFold(supernet.GroupName, UNUSED_VLANName) {
					// Unused Vlan defined in input json
					UNUSED_VLANID = supernet.IPv4.VlanID
				} else if strings.EqualFold(supernet.GroupName, NATIVE_VLANName) {
					// Native Vlan 99
					NATIVE_VLANID = supernet.IPv4.VlanID
				}
				// Assign the value
				vlanItem.GroupName = supernet.GroupName
				vlanItem.VlanName = supernet.IPv4.Name
				vlanItem.VlanID = supernet.IPv4.VlanID
				vlanItem.Cidr = supernet.IPv4.Cidr
				vlanItem.Subnet = supernet.IPv4.Subnet
				vlanItem.Mtu = JUMBOMTU

				if supernet.Shutdown {
					vlanItem.Shutdown = true
				}
				if supernet.IPv4.SwitchSVI {
					for _, ipv4 := range supernet.IPv4.Assignment {
						if strings.Contains(strings.ToUpper(ipv4.Name), strings.ToUpper(VIPGATEWAY)) {
							vlanItem.VIPAddress = ipv4.IP
							// Caculate Virtual Group ID as limited 1~255.
							VIDBase := 50
							vlanItem.VirtualGroupID = idx + VIDBase
						} else if strings.Contains(strings.ToUpper(ipv4.Name), strings.ToUpper(o.Switch.Type)) {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
						}
					}
				}
				if !(strings.EqualFold(supernet.GroupName, STORAGE) && strings.EqualFold(inputData.DeploymentPattern, SWITCHLESS)) {
					vlanList = append(vlanList, vlanItem)
				}
			} else {
				// L3 Interface Object
				for _, ipv4 := range supernet.IPv4.Assignment {
					if strings.EqualFold(ipv4.Name, o.Switch.Type) {
						// Assignment Type binds with Switch.Type
						l3IntfItem.IPAddress = ipv4.IP
					}
					// Update NbrIPAddress for IBGP Peer
					for _, switchObj := range o.SwitchPeer {
						if strings.EqualFold(ipv4.Name, switchObj.Type) {
							// Assignment Type binds with Switch.Type
							l3IntfItem.NbrIPAddress = ipv4.IP
						}
					}
					// Update NbrIPaddress for P2P_Border
					for _, switchObj := range o.SwitchUplink {
						if strings.EqualFold(ipv4.Name, switchObj.Type) {
							// Assignment Type binds with Switch.Type
							l3IntfItem.NbrIPAddress = ipv4.IP
						}
					}
				}
				if len(l3IntfItem.IPAddress) != 0 {
					l3IntfItem.Function = supernet.IPv4.Name
					l3IntfItem.Description = supernet.IPv4.Name
					l3IntfItem.Cidr = supernet.IPv4.Cidr
					l3IntfItem.Mtu = JUMBOMTU
					l3IntfItem.Subnet = supernet.IPv4.Subnet
					// Upper case the key Name
					l3IntfMap[strings.ToUpper(supernet.IPv4.Name)] = l3IntfItem
				}
			}
		}
	} else if strings.Contains(o.Switch.Type, BMC) {
		// BMC Switch only have Unused and BMC Vlan (no VIP)
		for _, supernet := range inputData.Supernets {
			if supernet.IPv4.VlanID != 0 {
				vlanItem := VlanType{}
				if strings.Contains(strings.ToUpper(supernet.GroupName), strings.ToUpper(BMC)) {
					// BMC Vlan
					BMC_VlanID = supernet.IPv4.VlanID
					vlanItem.GroupName = supernet.GroupName
					vlanItem.VlanName = supernet.IPv4.Name
					vlanItem.VlanID = supernet.IPv4.VlanID
					vlanItem.Cidr = supernet.IPv4.Cidr
					vlanItem.Subnet = supernet.IPv4.Subnet
					vlanItem.Mtu = JUMBOMTU
					vlanItem.Shutdown = false

					for _, ipv4 := range supernet.IPv4.Assignment {
						if strings.Contains(strings.ToUpper(ipv4.Name), strings.ToUpper(VIPGATEWAY)) {
							vlanItem.VIPAddress = ipv4.IP
						} else if strings.Contains(strings.ToUpper(ipv4.Name), strings.ToUpper(o.Switch.Type)) {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
						}
					}
					vlanList = append(vlanList, vlanItem)
				} else if strings.Contains(strings.ToUpper(supernet.GroupName), strings.ToUpper(UNUSED_VLANName)) {
					// Unused Vlan defined in input json
					vlanItem.VlanID = supernet.IPv4.VlanID
					vlanItem.GroupName = supernet.GroupName
					vlanItem.VlanName = supernet.IPv4.Name
					vlanItem.Mtu = JUMBOMTU
					vlanItem.Shutdown = true
					vlanList = append(vlanList, vlanItem)
				} else if strings.Contains(strings.ToUpper(supernet.GroupName), strings.ToUpper(NATIVE_VLANName)) {
					// Native Vlan 99
					vlanItem.VlanID = supernet.IPv4.VlanID
					vlanItem.GroupName = supernet.GroupName
					vlanItem.VlanName = supernet.IPv4.Name
					vlanItem.Mtu = JUMBOMTU
					vlanItem.Shutdown = false
					vlanList = append(vlanList, vlanItem)
				}
			}
		}
	}

	sort.Slice(vlanList, func(i, j int) bool {
		return vlanList[i].VlanID < vlanList[j].VlanID
	})

	o.Vlans = vlanList
	o.L3Interfaces = l3IntfMap
}
