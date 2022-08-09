package main

import "strings"

func newLoopbackObj() *LoopbackType {
	return &LoopbackType{}
}

func (o *OutputType) parseLoopbackObj(i *InterfaceFrameworkType) {
	for _, item := range i.Loopback {
		newlp := newLoopbackObj()
		newlp.Description = item.Description
		LoopbackName := strings.Replace(item.IPAddress, "TORX", o.Device.Type, -1)
		tmp := o.getSwitchMgmtIPbyName(LoopbackName)
		newlp.IPAddress = tmp
		o.Loopback = append(o.Loopback, *newlp)
	}
}
