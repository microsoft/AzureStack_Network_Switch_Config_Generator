package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (o *OutputType) parsePortchannelObj(i *InterfaceFrameworkType) {
	o.PortChannel = i.PortChannel
	for i, item := range o.PortChannel {
		if item.Type == IP && len(item.IPAddress) > 0 {
			PCIPName := strings.Replace(item.IPAddress, "TORX", TORX, -1)
			PCIPAddress := o.getIPbyName(PCIPName, SWITCH_MGMT)
			PCNbrIPName := strings.Replace(item.NbrIPAddress, "TORY", TORY, -1)
			PCNbrIPAddress := o.getIPbyName(PCNbrIPName, SWITCH_MGMT)
			o.PortChannel[i].IPAddress = PCIPAddress
			o.PortChannel[i].NbrIPAddress = PCNbrIPAddress
		} else if item.Type == TRUNK {
			tmpVlans := []string{}
			vlanList := strings.Split(item.VLANs, ",")
			for _, vlanName := range vlanList {
				if vlanName == CISCO_NATIVE_VLAN {
					tmpVlans = append(tmpVlans, CISCO_NATIVE_VLAN)
				} else {
					vlanID := o.getVLANIDbyName(vlanName)
					vlanIDStr := strconv.Itoa(vlanID)
					tmpVlans = append(tmpVlans, vlanIDStr)
				}
			}
			o.PortChannel[i].VLANs = strings.Join(tmpVlans, ",")
		}
		o.UpdatePortAttrbyPC(o.PortChannel[i])
	}
}

func (o *OutputType) UpdatePortAttrbyPC(p PortChannelType) {
	if len(p.Members) == 0 {
		log.Fatalf("PortChannel %d has 0 members\n", p.ID)
	}
	for _, m := range p.Members {
		for i, port := range o.Port {
			if m == port.Port {
				o.Port[i].Description = p.Description
				o.Port[i].PortType = p.Type
				o.Port[i].Shutdown = p.Shutdown
				o.Port[i].TagVlan = p.VLANs
				o.Port[i].UntagVlan = 0
				o.Port[i].Others = p.Others
				o.Port[i].Others[CHANNEL_GROUP] = fmt.Sprintf("%d", p.ID)
			}
		}
	}
}
