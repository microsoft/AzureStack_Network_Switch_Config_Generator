package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func (o *OutputType) ParseRouting(frameworkFolder string, inputData InputData) {
	routingJsonPath := fmt.Sprintf("%s/%s", frameworkFolder, BGPROUTINGJSON)
	routingJsonObj := parseRoutingJson(routingJsonPath)
	if inputData.SwitchUplink == "BGP" {
		o.ParseBGP(routingJsonObj.BGP)
		o.ParsePrefixList(routingJsonObj.PrefixList)
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
	newBGPObj.RouterID = o.getIPAddressByL3Intf(BGPObj.RouterID)

	// # Update BGP Advertised Network
	ipv4Networks := []string{}
	// ## All L3 Interfaces [P2P+Loopback]
	for _, l3intfItem := range o.L3Interfaces {
		ipv4Networks = append(ipv4Networks, l3intfItem.Subnet)
	}
	// ## Selected Vlan Subnet
	for _, networkName := range BGPObj.IPv4Network {
		subnetList := o.getSubnetByVlanGroupID(networkName)
		ipv4Networks = append(ipv4Networks, subnetList...)
	}
	newBGPObj.IPv4Network = ipv4Networks

	// # Update BGP Neighbor Object
	newBGPIPv4Nbrs := []IPv4NeighborType{}
	for _, ipv4NbrItem := range BGPObj.IPv4Neighbor {
		if ipv4NbrItem.SwitchRelation == "SwitchUplink" {
			for _, switchItem := range o.SwitchUplink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("To_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				newBGPIPv4NbrItem.NeighborIPAddress = o.getIPAddressByL3Intf(switchItem.Type)
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		} else if ipv4NbrItem.SwitchRelation == "SwitchPeer" {
			for _, switchItem := range o.SwitchPeer {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("To_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				newBGPIPv4NbrItem.NbrPassword = generateRandomString(16, 3, 3, 3)
				newBGPIPv4NbrItem.NeighborIPAddress = o.getIPAddressByL3Intf(ipv4NbrItem.NeighborIPAddress)
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		} else if ipv4NbrItem.SwitchRelation == "SwitchDownlink" {
			for _, switchItem := range o.SwitchDownlink {
				newBGPIPv4NbrItem := ipv4NbrItem
				newBGPIPv4NbrItem.Description = fmt.Sprintf("To_%s", switchItem.Type)
				newBGPIPv4NbrItem.NeighborAsn = switchItem.Asn
				newBGPIPv4NbrItem.NeighborIPAddress = o.getSubnetByVlanGroupID(ipv4NbrItem.NeighborIPAddress)[0]
				newBGPIPv4Nbrs = append(newBGPIPv4Nbrs, newBGPIPv4NbrItem)
			}
		}
	}
	newBGPObj.IPv4Neighbor = newBGPIPv4Nbrs
	// Assign back to final json Object
	o.Routing.BGP = newBGPObj
}

func (o *OutputType) ParsePrefixList(PrefixListJson []PrefixListType) {
	prefixListObj := []PrefixListType{}
	for _, prefixListTmp := range PrefixListJson {
		prefixName, prefixConfig := prefixListTmp.Name, prefixListTmp.Config
		prefixListObjItem := PrefixListType{}
		prefixListObjItem.Name = prefixName
		for _, configItemTmp := range prefixConfig {
			newConfigItem := configItemTmp
			subnetList := o.getSubnetByVlanGroupID(configItemTmp.Network)
			for _, subnet := range subnetList {
				newConfigItem.Network = subnet
				prefixListObjItem.Config = append(prefixListObjItem.Config, newConfigItem)
			}
		}
		prefixListObj = append(prefixListObj, prefixListObjItem)
	}
	o.Routing.PrefixList = prefixListObj
}

func (o *OutputType) getSubnetByVlanGroupID(groupID string) []string {
	sunbetList := []string{}
	if groupID == ANY {
		sunbetList = []string{"0.0.0.0/0"}
	} else {
		for _, vlanItem := range o.Vlans {
			if vlanItem.GroupID == groupID {
				sunbetList = append(sunbetList, vlanItem.Subnet)
			}
		}
	}
	return sunbetList
}

func (o *OutputType) getIPAddressByL3Intf(networkName string) string {
	var subnet string
	for key, l3inftItem := range o.L3Interfaces {
		if strings.Contains(key, networkName) {
			subnet = l3inftItem.IPAddress
		}
	}
	return subnet
}
