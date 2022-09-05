package main

import (
	"log"
	"strconv"
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

		if intf.PortType == PortType_BMC_MGMT {
			intf.IPAddress = replaceTORXName(intf.IPAddress, o.Device.Type)
			portObj.IPAddress = o.searchOOBIP(intf.IPAddress)
		} else if intf.PortType == PortType_IP {
			intf.IPAddress = replaceTORXName(intf.IPAddress, o.Device.Type)
			portObj.IPAddress = o.getSwitchMgmtIPbyName(intf.IPAddress)
		} else if intf.PortType == PortType_ACCESS {
			portObj.UntagVlan = o.updatePortUntagvlan(intf.UntagVlan)
		} else if intf.PortType == PortType_TRUNK {
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
		if segment.Group == PortType_BMC_MGMT {
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
		if segment.Group == PortType_INFRA_MGMT {
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

func (o *OutputType) updatePortTagvlan(TagVlanName []string) string {
	// if len(TagVlanName) < 1 {
	// 	log.Fatalln("Tag Vlan Attributes of this Trunk Port is invalid.")
	// }
	res := []string{}
	for _, segment := range *o.Supernets {
		for _, vlanName := range TagVlanName {
			if segment.VlanID != 0 {
				if segment.Group == vlanName {
					res = append(res, strconv.Itoa(segment.VlanID))
				}
			}
		}
	}
	return strings.Join(res, ",")
}

func (i *InterfaceFrameworkType) parseInterfaceName(Speed int) string {
	for _, intfName := range i.InterfaceName {
		if intfName.Speed == Speed {
			return intfName.Name
		}
	}
	return ""
}
