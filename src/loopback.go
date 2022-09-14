package main

import (
	"strings"
)

func newLoopbackObj() *LoopbackType {
	return &LoopbackType{}
}

func (o *OutputType) parseLoopbackObj(i *InterfaceFrameworkType) {
	for _, item := range i.Loopback {
		newlp := newLoopbackObj()
		newlp.Name = item.Name
		LoopbackName := strings.Replace(item.IPAddress, "TORX", o.Device.Type, -1)
		tmp := o.getIPbyName(LoopbackName, SWITCH_MGMT)
		newlp.IPAddress = tmp
		o.Loopback = append(o.Loopback, *newlp)
	}
}
