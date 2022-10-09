package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func newRoutingObj() *RoutingType {
	return &RoutingType{}
}

func parseRoutingJSON(routingFrameJson string) *RoutingType {
	routingFrameworkObj := newRoutingObj()
	bytes, err := ioutil.ReadFile(routingFrameJson)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, routingFrameworkObj)
	if err != nil {
		log.Fatalln(err)
	}
	return routingFrameworkObj
}

func (o *OutputType) parseRoutingFramework(frameworkPath, deviceType string, inputJsonObj *InputType) {
	routingFrameJson := fmt.Sprintf("%s/%s_%s.%s", frameworkPath, deviceType, ROUTING, JSON)
	routingFrameworkObj := parseRoutingJSON(strings.ToLower(routingFrameJson))
	// Use StaticRouting attribute to update StaticRoute or BGPRoute
	if o.Device.StaticRouting {
		// Static Routing
		routingFrameworkObj.updateStaticRoutingPolicy(o)
		routingFrameworkObj.updateStaticNetwork(o)
	} else {
		// BGP Routing
		routingFrameworkObj.Bgp.BGPAsn = strconv.Itoa(o.Device.Asn)
		routingFrameworkObj.updateBgpNetwork(o)
		routingFrameworkObj.updateBGPRoutingPolicy(o)
		routerIDName := strings.Replace(routingFrameworkObj.Bgp.RouterID, "TORX", TORX, -1)
		RouterIDIPAddress := o.getIPbyName(routerIDName, SWITCH_MGMT)
		routingFrameworkObj.Bgp.RouterID = strings.Split(RouterIDIPAddress, "/")[0]
		routingFrameworkObj.updateBgpNeighbor(o, inputJsonObj)
	}
	o.Routing = routingFrameworkObj
}

func (r *RoutingType) updateBgpNetwork(outputObj *OutputType) {
	for index, netname := range r.Bgp.IPv4Network {
		r.Bgp.IPv4Network[index] = outputObj.getSupernetIPbyName(netname)
	}
}

func (r *RoutingType) updateBGPRoutingPolicy(outputObj *OutputType) {
	newPrefixList := []PrefixListType{}
	// Update PrefexList
	for i, item := range r.RoutingPolicy.PrefixList {
		for _, prefixListName := range r.Bgp.PrefixListName {
			if prefixListName == item.Name {
				for j, config := range item.Config {
					// Update Supernet Name to Supernet IP
					supernetIP := outputObj.getSupernetIPbyName(config.Supernet)
					r.RoutingPolicy.PrefixList[i].Config[j].Supernet = supernetIP
				}
				//Append the validate item to newPrefixList
				newPrefixList = append(newPrefixList, item)
			}
		}
	}
	r.RoutingPolicy.PrefixList = newPrefixList
	r.RoutingPolicy.RouteMap = nil
}

func (r *RoutingType) updateBgpNeighbor(outputObj *OutputType, inputJsonObj *InputType) {
	for k, v := range r.Bgp.IPv4Neighbor {
		nbrASNName := strings.Replace(v.NeighborAsn, "TORY", TORY, -1)
		nbrAsn, err := inputJsonObj.getBgpASN(nbrASNName)
		if err != nil {
			log.Fatalln(err)
		}
		r.Bgp.IPv4Neighbor[k].NeighborAsn = nbrAsn
		nbrIPAddressName := strings.Replace(v.NeighborIPAddress, "TORX", TORX, -1)
		nbrIPAddressName = strings.Replace(nbrIPAddressName, "TORY", TORY, -1)
		IPv4IPNet := outputObj.getIPbyName(nbrIPAddressName, SWITCH_MGMT)
		r.Bgp.IPv4Neighbor[k].NeighborIPAddress = strings.Split(IPv4IPNet, "/")[0]
		r.Bgp.IPv4Neighbor[k].Description = nbrIPAddressName

		updateSourceName := replaceTORXName(v.UpdateSource, outputObj.Device.Type)
		r.Bgp.IPv4Neighbor[k].UpdateSource = outputObj.getIPbyName(updateSourceName, SWITCH_MGMT)
	}
}

func (r *RoutingType) updateStaticNetwork(outputObj *OutputType) {
	tmp := []StaticNetworkType{}
	for _, staticItem := range r.Static.Network {
		// Replace template name with right TOR number.
		routeName := strings.Replace(staticItem.Name, "TORX", TORX, -1)
		if len(staticItem.NextHop) != 0 {
			// Update null 0 static route or Get BMCmgmt VIP
			if strings.Contains(staticItem.NextHop, DeviceType_BMC) {
				staticItem.NextHop = outputObj.getIPbyName(staticItem.NextHop, BMC_MGMT)
			}
			tmp = append(tmp, StaticNetworkType{
				DstIPAddress: outputObj.getSupernetIPbyName(staticItem.DstIPAddress),
				NextHop:      staticItem.NextHop,
				Name:         routeName,
			})
		} else if len(staticItem.DstIPAddress) != 0 {
			// update default route to border
			nextHop := outputObj.getIPbyName(routeName, SWITCH_MGMT)
			tmp = append(tmp, StaticNetworkType{
				DstIPAddress: outputObj.getSupernetIPbyName(staticItem.DstIPAddress),
				NextHop:      strings.Split(nextHop, "/")[0],
				Name:         routeName,
			})
		}
	}
	r.Static.Network = tmp
}

func (o *OutputType) getSupernetIPbyName(SupernetName string) string {
	for _, segment := range *o.Supernets {
		if segment.Name == SupernetName {
			return segment.Subnet
		}
	}
	return ""
}

func (r *RoutingType) updateStaticRoutingPolicy(outputObj *OutputType) {
	newPrefixList := []PrefixListType{}
	// Update PrefexList
	for i, item := range r.RoutingPolicy.PrefixList {
		for _, prefixListName := range r.Static.PrefixListName {
			if prefixListName == item.Name {
				for j, config := range item.Config {
					// Update Supernet Name to Supernet IP
					supernetIP := outputObj.getSupernetIPbyName(config.Supernet)
					r.RoutingPolicy.PrefixList[i].Config[j].Supernet = supernetIP
				}
				//Append the validate item to newPrefixList
				newPrefixList = append(newPrefixList, item)
			}
		}
	}
	r.RoutingPolicy.PrefixList = newPrefixList
}
