package main

import (
	"strings"
)

func newVlanObj() *VlanType {
	return &VlanType{}
}

func (o *OutputType) parseVlanObj(i *InterfaceFrameworkType) {

	for _, segment := range *o.Supernets {
		if segment.VlanID > 0 {
			vlanObj := newVlanObj()
			vlanObj.VlanName = segment.Name
			vlanObj.Shutdown = segment.Shutdown
			vlanObj.VlanID = segment.VlanID
			vlanObj.Group = segment.Group
			for _, vlanItem := range i.Vlan {
				if vlanObj.Group == vlanItem.Group {
					vlanObj.Mtu = vlanItem.Mtu
					IPAssignmentName := replaceTORXName(vlanItem.IPAssignment, o.Device.Type)
					for _, ipAssign := range segment.IPAssignment {
						if strings.Contains(ipAssign.Name, IPAssignmentName) {
							vlanObj.IPAddress = ipAssign.IPAddress
						}
					}
				}
			}
			o.Vlan = append(o.Vlan, *vlanObj)
		}
	}
}
