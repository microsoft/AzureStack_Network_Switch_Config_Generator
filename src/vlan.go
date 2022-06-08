package main

import (
	"strings"
)

func newVlanObj() *VlanType {
	return &VlanType{}
}

func (o *OutputJsonType) parseVlanFramework(i *InterfaceFrameworkType) {
	for _, vlanItem := range i.Vlan {
		vlanObj := newVlanObj()
		vlanObj.VlanName = strings.Replace(vlanItem.Name, "TORX", o.Device.Type, -1)
		vlanObj.Mtu = vlanItem.Mtu
		vlanObj.IPAddress = o.updateVlanIPAddr(vlanItem.Name, vlanItem.IPAssignment)
		vlanObj.updateVlanFromNetwork(o, vlanItem.Name)
		o.Vlan = append(o.Vlan, *vlanObj)
	}
}

func (o *OutputJsonType) updateVlanIPAddr(VlanName, IPAssignment string) string {
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

func (v *VlanType) updateVlanFromNetwork(o *OutputJsonType, VlanName string) {
	for _, segment := range *o.Network {
		if segment.Name == VlanName {
			v.Shutdown = segment.Shutdown
			v.Type = segment.Type
			v.VlanID = segment.VlanID
		}
	}
}
