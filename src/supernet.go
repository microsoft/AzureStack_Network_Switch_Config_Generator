package main

import (
	"strings"
)

func (o *OutputType) UpdateVlanAndL3Intf(inputData InputData) {
	vlanList := []VlanType{}
	l3IntfMap := map[string]L3IntfType{}
	if strings.Contains(o.Switch.Type, TOR) {
		// TOR Switch with all matched vlans
		for _, supernet := range inputData.Supernets {
			vlanItem := VlanType{}
			l3IntfItem := L3IntfType{}
			if supernet.IPv4.VlanID != 0 {
				if strings.Contains(supernet.GroupName, BMC) {
					BMC_VlanID = supernet.IPv4.VlanID
				} else if strings.Contains(supernet.GroupName, Infra_GroupID) {
					Infra_VlanID = supernet.IPv4.VlanID
				}
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
						if ipv4.Name == VIPGATEWAY {
							vlanItem.VIPAddress = ipv4.IP
						} else if ipv4.Name == o.Switch.Type {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
						}
					}
				}
				vlanList = append(vlanList, vlanItem)
			} else {
				// L3 Interface Object
				for _, ipv4 := range supernet.IPv4.Assignment {
					if ipv4.Name == o.Switch.Type {
						// Assignment Type binds with Switch.Type
						l3IntfItem.IPAddress = ipv4.IP
					}
					// Update NbrIPAddress for IBGP Peer
					for _, switchObj := range o.SwitchPeer {
						if ipv4.Name == switchObj.Type {
							// Assignment Type binds with Switch.Type
							l3IntfItem.NbrIPAddress = ipv4.IP
						}
					}
					// Update NbrIPaddress for P2P_Border
					for _, switchObj := range o.SwitchUplink {
						if ipv4.Name == switchObj.Type {
							// Assignment Type binds with Switch.Type
							l3IntfItem.NbrIPAddress = ipv4.IP
						}
					}
				}
				if len(l3IntfItem.IPAddress) != 0 {
					l3IntfItem.Function = supernet.IPv4.NetworkType
					l3IntfItem.Description = supernet.IPv4.Name
					l3IntfItem.Cidr = supernet.IPv4.Cidr
					l3IntfItem.Mtu = JUMBOMTU
					l3IntfItem.Subnet = supernet.IPv4.Subnet
					l3IntfMap[supernet.IPv4.Name] = l3IntfItem
				}
			}
		}
	} else if strings.Contains(o.Switch.Type, BMC) {
		// BMC Switch only have Unused and BMC Vlan (no VIP)
		for _, supernet := range inputData.Supernets {
			vlanItem := VlanType{}
			if supernet.GroupName == UNUSED || supernet.GroupName == BMC {
				vlanItem.VlanName = supernet.IPv4.Name
				vlanItem.VlanID = supernet.IPv4.VlanID
				vlanItem.GroupName = supernet.GroupName
				vlanItem.Mtu = DefaultMTU
				if supernet.Shutdown {
					vlanItem.Shutdown = true
				}
				if supernet.IPv4.SwitchSVI {
					for _, ipv4 := range supernet.IPv4.Assignment {
						if strings.Contains(ipv4.Name, o.Switch.Type) {
							// Assignment Type binds with Switch.Type
							vlanItem.IPAddress = ipv4.IP
							vlanItem.Cidr = supernet.IPv4.Cidr
							vlanItem.Subnet = supernet.IPv4.Subnet
							vlanItem.VIPAddress = supernet.IPv4.Gateway
						}
					}
				}
				vlanList = append(vlanList, vlanItem)
			}
		}
	}

	o.Vlans = vlanList
	o.L3Interfaces = l3IntfMap
}
