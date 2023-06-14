package main

import (
	"fmt"
	"strings"
)

func (o *OutputType) UpdatePortChannel(inputData InputData) {
	portChannelMap := map[string]PortChannelType{}
	PO_TOR_BMC := PortChannelType{
		Description:   TOR_BMC,
		Function:      TOR_BMC,
		UntagVlan:     NATIVE_VLANID,
		TagVlans:      BMC_VlanID,
		PortChannelID: POID_TOR_BMC,
		VPC:           POID_TOR_BMC,
		Shutdown:      false,
	}
	if strings.Contains(o.Switch.Type, TOR) {
		// TOR Switch has all PortChannels
		PO_MLAG_PEER := PortChannelType{
			Description:   MLAG_PEER,
			Function:      MLAG_PEER,
			UntagVlan:     NATIVE_VLANID,
			PortChannelID: POID_MLAG_PEER,
			VPC:           "peer-link",
			Shutdown:      false,
		}
		var P2P_IBGP_IP string
		for _, l3IntfItem := range o.L3Interfaces {
			if strings.EqualFold(l3IntfItem.Function, P2P_IBGP) {
				P2P_IBGP_IP = fmt.Sprintf("%s/%d", l3IntfItem.IPAddress, l3IntfItem.Cidr)
			}
		}

		PO_P2P_IBGP := PortChannelType{
			Description:   P2P_IBGP,
			Function:      P2P_IBGP,
			IPAddress:     P2P_IBGP_IP,
			PortChannelID: POID_P2P_IBGP,
			Shutdown:      false,
		}
		portChannelMap[MLAG_PEER] = PO_MLAG_PEER
		portChannelMap[P2P_IBGP] = PO_P2P_IBGP
		portChannelMap[TOR_BMC] = PO_TOR_BMC

	} else if strings.Contains(o.Switch.Type, BMC) {
		// BMC Switch only have PO_TOR_BMC
		portChannelMap[TOR_BMC] = PO_TOR_BMC
	}

	// if strings.Contains(o.Switch.Type, TOR) && len(o.SwitchBMC) == 0 {
	// 	delete(portChannelMap, TOR_BMC)
	// }

	o.PortChannel = portChannelMap
}
