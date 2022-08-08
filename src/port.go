package main

import (
	"log"
	"strings"
)

func newPortObj() *PortType {
	return &PortType{}
}

func (o *OutputType) parseInBandPortFramework(i *InterfaceFrameworkType) {
	for _, intf := range i.Port {
		portObj := newPortObj()
		portObj.Port = intf.Port
		intf.Description = strings.Replace(intf.Description, "TORX", o.Device.Type, -1)
		portObj.Description = intf.Description
		portObj.Mtu = intf.Mtu
		portObj.PortType = intf.PortType
		portObj.Shutdown = intf.Shutdown
		portObj.Others = intf.Others
		if intf.Speed != 0 {
			portObj.PortName = i.parseInterfaceName(intf.Speed)
		} else {
			portObj.PortName = intf.Port
		}

		if intf.PortType == "OOB" {
			intf.IPAddress = replaceTORXName(intf.IPAddress, o.Device.Type)
			portObj.IPAddress = o.searchOOBIP(intf.IPAddress)
		} else if intf.PortType == "IP" {
			intf.IPAddress = replaceTORXName(intf.IPAddress, o.Device.Type)
			portObj.IPAddress = o.getSwitchMgmtIPbyName(intf.IPAddress)
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

func (o *OutputType) searchOOBIP(IPAddressName string) string {
	for _, segment := range *o.Supernets {
		if segment.Group == "OOB" {
			for _, ipAssign := range segment.IPAssignment {
				if strings.Contains(ipAssign.Name, IPAddressName) {
					return ipAssign.IPAddress
				}
			}
		}
	}
	return ""
}

func (o *OutputType) getSwitchMgmtIPbyName(SwitchMgmtName string) string {
	for _, segment := range *o.Supernets {
		if segment.Group == "SwitchMgmt" {
			for _, ipAssign := range segment.IPAssignment {
				if strings.Contains(ipAssign.Name, SwitchMgmtName) {
					return ipAssign.IPAddress
				}
			}
		}
	}
	return ""
}

func (o *OutputType) updatePortUntagvlan(UntagVlanName string) int {
	for _, segment := range *o.Supernets {
		if segment.VlanID != 0 {
			if segment.Name == UntagVlanName {
				return segment.VlanID
			}
		}
	}
	return 0
}

func (o *OutputType) updatePortTagvlan(TagVlanName []string) []int {
	if len(TagVlanName) < 1 {
		log.Fatalln("Tag Vlan Attributes of this Trunk Port is invalid.")
	}
	ret := make([]int, len(TagVlanName))
	for _, segment := range *o.Supernets {
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
