package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func (o *OutputType) ParseRouting(frameworkFolder string, inputData InputData) {
	if strings.Contains(strings.ToUpper(o.Switch.Type), TOR) {
		if strings.ToUpper(inputData.SwitchUplink) == BGP {
			routingJsonPath := fmt.Sprintf("%s/%s.json", frameworkFolder, strings.ToLower(BGP))
			routingJsonObj := parseRoutingJson(routingJsonPath)
			o.ParseBGP(routingJsonObj.BGP)
			o.ParsePrefixList(routingJsonObj.PrefixList)
		} else if strings.ToUpper(inputData.SwitchUplink) == STATIC {
			routingJsonPath := fmt.Sprintf("%s/%s.json", frameworkFolder, strings.ToLower(STATIC))
			routingJsonObj := parseRoutingJson(routingJsonPath)
			o.ParseBGP(routingJsonObj.BGP)
			o.ParsePrefixList(routingJsonObj.PrefixList)
			o.ParseStatic(routingJsonObj.Static)
		}
	} else if strings.Contains(strings.ToUpper(o.Switch.Type), BMC) {
		routingJsonPath := fmt.Sprintf("%s/%s.json", frameworkFolder, strings.ToLower(STATIC))
		routingJsonObj := parseRoutingJson(routingJsonPath)
		o.ParseStatic(routingJsonObj.Static)
	}
}

func parseRoutingJson(routingJsonPath string) *RoutingType {
	routingJsonObj := &RoutingType{}
	bytes, err := ioutil.ReadFile(routingJsonPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, routingJsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	return routingJsonObj
}

func (o *OutputType) ParseBGP(BGPObj BGPType) {
	newBGPObj := BGPObj
	// # Update BGPASN
	newBGPObj.BGPAsn = o.Switch.Asn
	newBGPObj.RouterID = o.getL3IntfObjByName(BGPObj.RouterID).IPAddress

	// # Update BGP Advertised Network
	ipv4Networks := []string{}
	// ## All L3 Interfaces [P2P+Loopback]
	for _, l3intfItem := range o.L3Interfaces {
		ipv4Networks = append(ipv4Networks, l3intfItem.Subnet)
	}
	// ## Selected Vlan Subnet
	for _, networkName := range BGPObj.IPv4Network {
		subnetMap := o.getSubnetByVlanGroupID(networkName)
		for _, subnet := range subnetMap {
			ipv4Networks = append(ipv4Networks, subnet)
		}
	}
	newBGPObj.IPv4Network = ipv4Networks

	// # Update BGP Neighbor Object
	newBGPIPv4Nbrs := []IPv4NeighborType{}
	for _, ipv4NbrItem := range BGPObj.IPv4Neighbor {
		if ipv4NbrItem.SwitchRelation == SWITCHUPLINK {
			for _, switchItem := range o.SwitchUplink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("TO_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				newBGPIPv4NbrItem.NeighborIPAddress = o.getL3IntfObjByName(switchItem.Type).NbrIPAddress
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		} else if ipv4NbrItem.SwitchRelation == SWITCHPEER {
			for _, switchItem := range o.SwitchPeer {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("TO_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				newBGPIPv4NbrItem.NbrPassword = generateRandomString(16, 3, 3, 3)
				newBGPIPv4NbrItem.NeighborIPAddress = o.getL3IntfObjByName(ipv4NbrItem.NeighborIPAddress).NbrIPAddress
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		} else if ipv4NbrItem.SwitchRelation == SWITCHDOWNLINK {
			for _, switchItem := range o.SwitchDownlink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("TO_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				for _, subnetValue := range o.getSubnetByVlanGroupID(ipv4NbrItem.NeighborIPAddress) {
					newBGPIPv4NbrItem.NeighborIPAddress = subnetValue
				}
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		}
	}
	newBGPObj.IPv4Neighbor = newBGPIPv4Nbrs
	// Assign back to final json Object
	o.Routing.BGP = newBGPObj

	if o.Switch.Type == BMC {
		o.Routing.BGP = BGPType{}
	}
}

func (o *OutputType) ParseStatic(staticObj []StaticType) {
	newStaticObj := []StaticType{}
	for _, staticItem := range staticObj {
		newStaticItem := staticItem
		if staticItem.Network != ANY {
			subnetMap := o.getSubnetByVlanGroupID(staticItem.Network)
			for networkName, network := range subnetMap {
				newStaticItem.Network = network
				newStaticItem.Name = networkName
				newStaticObj = append(newStaticObj, newStaticItem)
			}
		} else {
			if staticItem.NextHop == SWITCHUPLINK {
				for _, switchItem := range o.SwitchUplink {
					l3IntfObj := o.getL3IntfObjByName(switchItem.Type)
					newStaticItem.Network = ANYNETWORK
					newStaticItem.NextHop = l3IntfObj.IPAddress
					newStaticItem.Name = l3IntfObj.Function
					newStaticObj = append(newStaticObj, newStaticItem)
				}
			} else {
				subnetMap := o.getVIPByVlanGroupID(staticItem.NextHop)
				for _, network := range subnetMap {
					newStaticItem.Network = ANYNETWORK
					newStaticItem.NextHop = network
					newStaticObj = append(newStaticObj, newStaticItem)
				}
			}
		}
	}
	// Assign back to final json Object
	o.Routing.Static = newStaticObj
}

func (o *OutputType) ParsePrefixList(PrefixListObj []PrefixListType) {
	prefixListObj := []PrefixListType{}
	for _, prefixListTmp := range PrefixListObj {
		prefixName, prefixConfig := prefixListTmp.Name, prefixListTmp.Config
		prefixListObjItem := PrefixListType{}
		prefixListObjItem.Name = prefixName
		for _, configItemTmp := range prefixConfig {
			newConfigItem := configItemTmp
			subnetMap := o.getSubnetByVlanGroupID(configItemTmp.Network)
			for _, subnet := range subnetMap {
				newConfigItem.Network = subnet
				prefixListObjItem.Config = append(prefixListObjItem.Config, newConfigItem)
			}
		}
		prefixListObj = append(prefixListObj, prefixListObjItem)
	}
	o.Routing.PrefixList = prefixListObj
}

func (o *OutputType) getSubnetByVlanGroupID(groupID string) map[string]string {
	sunbetMap := map[string]string{}
	if groupID == ANY {
		sunbetMap = map[string]string{ANY: ANYNETWORK}
	} else {
		for _, vlanItem := range o.Vlans {
			if vlanItem.GroupName == groupID {
				sunbetMap[vlanItem.VlanName] = vlanItem.Subnet
			}
		}
	}
	return sunbetMap
}

func (o *OutputType) getVIPByVlanGroupID(groupID string) map[string]string {
	sunbetMap := map[string]string{}
	if groupID == ANY {
		sunbetMap = map[string]string{ANY: ANYNETWORK}
	} else {
		for _, vlanItem := range o.Vlans {
			if vlanItem.GroupName == groupID {
				sunbetMap[vlanItem.VlanName] = vlanItem.VIPAddress
			}
		}
	}
	return sunbetMap
}

func (o *OutputType) getL3IntfObjByName(networkName string) L3IntfType {
	var L3IntfObj L3IntfType
	for key, l3inftItem := range o.L3Interfaces {
		if strings.Contains(strings.ToUpper(key), strings.ToUpper(networkName)) {
			L3IntfObj = l3inftItem
		}
	}
	return L3IntfObj
}
