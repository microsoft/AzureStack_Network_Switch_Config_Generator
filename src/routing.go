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
		_, vlanSubnetList := o.getSubnetByVlanGroupName(networkName)
		ipv4Networks = append(ipv4Networks, vlanSubnetList...)
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
				newBGPIPv4NbrItem.NbrPassword = o.GlobalSetting.Password
				newBGPIPv4NbrItem.NeighborIPAddress = o.getL3IntfObjByName(ipv4NbrItem.NeighborIPAddress).NbrIPAddress
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		} else if ipv4NbrItem.SwitchRelation == SWITCHDOWNLINK {
			for _, switchItem := range o.SwitchDownlink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("TO_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				_, vlanSubnetList := o.getSubnetByVlanGroupName(ipv4NbrItem.NeighborIPAddress)
				for _, subnetValue := range vlanSubnetList {
					newBGPIPv4NbrItem.NeighborIPAddress = subnetValue
				}
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		}
	}
	// Dell Template SDN for Mux BGP
	newTemplateNeigbor := []IPv4NeighborType{}
	for _, ipv4NbrItem := range BGPObj.TemplateNeigbor {
		if ipv4NbrItem.SwitchRelation == SWITCHDOWNLINK {
			for _, switchItem := range o.SwitchDownlink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("TO_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				_, vlanSubnetList := o.getSubnetByVlanGroupName(ipv4NbrItem.NeighborIPAddress)
				for _, subnetValue := range vlanSubnetList {
					newBGPIPv4NbrItem.NeighborIPAddress = subnetValue
				}
				newTemplateNeigbor = append(newTemplateNeigbor, newBGPIPv4NbrItem)
			}
		}
	}
	newBGPObj.IPv4Neighbor = newBGPIPv4Nbrs
	newBGPObj.TemplateNeigbor = newTemplateNeigbor
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
			vlanNameList, vlanSubnetList := o.getSubnetByVlanGroupName(staticItem.Network)
			for index, subnet := range vlanSubnetList {
				newStaticItem.Network = subnet
				newStaticItem.Name = vlanNameList[index]
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
				subnetMap := o.getVIPByVlanGroupName(staticItem.NextHop)
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
			baseIndex := configItemTmp.Idx
			_, vlanSubnetList := o.getSubnetByVlanGroupName(configItemTmp.Network)
			for idx, subnet := range vlanSubnetList {
				newConfigItem.Network = subnet
				// Define the interval idx 5 for each ACL
				newConfigItem.Idx = baseIndex + idx*5
				prefixListObjItem.Config = append(prefixListObjItem.Config, newConfigItem)
			}
		}
		prefixListObj = append(prefixListObj, prefixListObjItem)
	}
	o.Routing.PrefixList = prefixListObj
}

func (o *OutputType) getSubnetByVlanGroupName(GroupName string) ([]string, []string) {
	vlanNameList, vlanSubnetList := []string{}, []string{}
	if strings.EqualFold(GroupName, ANY) {
		vlanNameList, vlanSubnetList = []string{ANY}, []string{ANYNETWORK}
	} else {
		for _, vlanItem := range o.Vlans {
			if strings.EqualFold(vlanItem.GroupName, GroupName) {
				vlanNameList = append(vlanNameList, vlanItem.VlanName)
				vlanSubnetList = append(vlanSubnetList, vlanItem.Subnet)
			}
		}
	}
	return vlanNameList, vlanSubnetList
}

func (o *OutputType) getVIPByVlanGroupName(GroupName string) map[int]string {
	sunbetMap := map[int]string{}
	GroupName = strings.ToUpper(GroupName)
	if GroupName == ANY {
		sunbetMap = map[int]string{0: ANYNETWORK}
	} else {
		for _, vlanItem := range o.Vlans {
			if vlanItem.GroupName == GroupName {
				sunbetMap[vlanItem.VlanID] = vlanItem.VIPAddress
			}
		}
	}

	return sunbetMap
}

func (o *OutputType) getIPByVlanGroupName(GroupName string) map[int]string {
	sunbetMap := map[int]string{}
	GroupName = strings.ToUpper(GroupName)

	for _, vlanItem := range o.Vlans {
		if vlanItem.GroupName == GroupName {
			sunbetMap[vlanItem.VlanID] = vlanItem.IPAddress
		}
	}
	return sunbetMap
}

func (o *OutputType) getL3IntfObjByName(networkName string) L3IntfType {
	var L3IntfObj L3IntfType
	networkName = strings.ToUpper(networkName)
	for key, l3inftItem := range o.L3Interfaces {
		if strings.Contains(strings.ToUpper(key), networkName) {
			L3IntfObj = l3inftItem
		}
	}
	return L3IntfObj
}
