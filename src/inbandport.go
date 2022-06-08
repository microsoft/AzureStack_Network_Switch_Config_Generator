package main

import (
	"log"
	"strings"
)

func newPortObj() *PortType {
	return &PortType{}
}

func (o *OutputJsonType) parseInBandPortFramework(i *InterfaceFrameworkType) {
	for _, intf := range i.InBandPort {
		portObj := newPortObj()
		portObj.Port = intf.Port
		intf.Description = strings.Replace(intf.Description, "TORX", o.Device.Type, -1)
		portObj.Description = intf.Description
		portObj.Mtu = intf.Mtu
		portObj.PortType = intf.PortType
		portObj.Shutdown = intf.Shutdown
		portObj.PortName = i.parseInterfaceName(intf.Speed)
		if intf.PortType == "IP" {
			intf.IPAddress = strings.Replace(intf.IPAddress, "TORX", o.Device.Type, -1)
			portObj.IPAddress = o.searchSwitchMgmtIP(intf.IPAddress)
		} else if intf.PortType == "Access" {
			portObj.UntagVlan = o.updatePortUntagvlan(intf.UntagVlan)
		} else if intf.PortType == "Trunk" {
			portObj.UntagVlan = o.updatePortUntagvlan(intf.UntagVlan)
			portObj.TagVlan = o.updatePortTagvlan(intf.TagVlan)
		} else {
			log.Fatalf("Port Type: %s is invalid\n", intf.PortType)
		}
		o.Port = append(o.Port, *portObj)
	}
}

func (o *OutputJsonType) searchSwitchMgmtIP(IPAddressName string) string {
	for _, segment := range *o.Network {
		if segment.Name == "SwitchMgmt" {
			for _, ipAssign := range segment.IPAssignment {
				if strings.Contains(ipAssign.Name, IPAddressName) {
					return ipAssign.IPAddress
				}
			}
		}
	}
	return ""
}

func (o *OutputJsonType) updatePortUntagvlan(UntagVlanName string) int {
	for _, segment := range *o.Network {
		if segment.VlanID != 0 {
			if segment.Name == UntagVlanName {
				return segment.VlanID
			}
		}
	}
	return 0
}

func (o *OutputJsonType) updatePortTagvlan(TagVlanName []string) []int {
	if len(TagVlanName) < 1 {
		log.Fatalln("Tag Vlan Attributes of this Trunk Port is invalid.")
	}
	ret := make([]int, len(TagVlanName))
	for _, segment := range *o.Network {
		for index, vlanName := range TagVlanName {
			if segment.VlanID != 0 {
				if segment.Name == vlanName {
					ret[index] = segment.VlanID
				}
			}
		}
	}
	return ret
}

func (i *InterfaceFrameworkType) parseInterfaceName(Speed int) string {
	for _, intfName := range i.InterfaceName {
		if intfName.Speed == Speed {
			return intfName.Name
		}
	}
	return ""
}
