package main

import (
	"strings"
)

func newVlanObj() *VlanType {
	return &VlanType{}
}

func (o *OutputType) parseVlanFramework(i *InterfaceFrameworkType) {
	for _, vlanItem := range i.Vlan {
		vlanObj := newVlanObj()
		vlanObj.VlanName = vlanItem.Name
		vlanObj.Mtu = vlanItem.Mtu
		IPAssignment := replaceTORXName(vlanItem.IPAssignment, o.Device.Type)
		vlanObj.IPAddress = o.updateVlanIPAddr(vlanItem.Name, IPAssignment)
		vlanObj.updateVlanFromNetwork(o, vlanItem.Name)
		o.Vlan = append(o.Vlan, *vlanObj)
	}
}

func (o *OutputType) updateVlanIPAddr(VlanName, IPAssignment string) string {
	for _, segment := range *o.Network {
		if segment.Name == VlanName {
			for _, ipAssign := range segment.IPAssignment {
				if strings.Contains(ipAssign.Name, IPAssignment) {
					return ipAssign.IPAddress
				}
			}
		}
	}
	return ""
}

func (v *VlanType) updateVlanFromNetwork(o *OutputType, VlanName string) {
	for _, segment := range *o.Network {
		if segment.Name == VlanName {
			v.Shutdown = segment.Shutdown
			v.Type = segment.Type
			v.VlanID = segment.VlanID
		}
	}
}
