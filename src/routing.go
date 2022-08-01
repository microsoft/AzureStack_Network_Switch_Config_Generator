package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (o *OutputType) parseRoutingFramework(frameworkPath string, inputJsonObj *InputType) {
	routingFrameJson := fmt.Sprintf("%s/routing.json", frameworkPath)
	routingFrameworkObj := parseRoutingJSON(routingFrameJson)
	// Use StaticRouting attribute to update StaticRoute or BGPRoute
	if o.Device.StaticRouting {
		// Static Routing
		routingFrameworkObj.updateStaticPrefixList(o)
		routingFrameworkObj.updateStaticNetwork(o)
	} else {
		// BGP Routing
		routingFrameworkObj.Bgp.BGPAsn = o.Device.Asn
		routingFrameworkObj.updateBgpNetwork(o)
		routerIDName := strings.Replace(routingFrameworkObj.Bgp.RouterID, "TORX", o.Device.Type, -1)
		routingFrameworkObj.Bgp.RouterID = o.getSwitchMgmtIPbyName(routerIDName)
		routingFrameworkObj.updateBgpNeighbor(o, inputJsonObj)
	}
	o.Routing = routingFrameworkObj
}

func (r *RoutingType) updateBgpNetwork(outputObj *OutputType) {
	for index, netname := range r.Bgp.IPv4Network {
		r.Bgp.IPv4Network[index] = outputObj.getSupernetIPbyName(netname)
	}
}

func (r *RoutingType) updateBgpNeighbor(outputObj *OutputType, inputJsonObj *InputType) {
	for k, v := range r.Bgp.IPv4Neighbor {
		nbrAsn, err := inputJsonObj.getBgpASN(v.NeighborAsn)
		if err != nil {
			log.Fatalln(err)
		}
		r.Bgp.IPv4Neighbor[k].NeighborAsn = nbrAsn
		nbrIPAddressName := replaceTORXName(v.NeighborIPAddress, outputObj.Device.Type)
		IPv4IPNet := outputObj.getSwitchMgmtIPbyName(nbrIPAddressName)
		r.Bgp.IPv4Neighbor[k].NeighborIPAddress = strings.Split(IPv4IPNet, "/")[0]
		r.Bgp.IPv4Neighbor[k].Description = nbrIPAddressName

		updateSourceName := replaceTORXName(v.UpdateSource, outputObj.Device.Type)
		r.Bgp.IPv4Neighbor[k].UpdateSource = outputObj.getSwitchMgmtIPbyName(updateSourceName)
	}
}

func (r *RoutingType) updateStaticNetwork(outputObj *OutputType) {
	tmp := []StaticNetworkType{}
	for _, staticItem := range r.Static.Network {
		// Replace template name with right TOR number.
		routeName := strings.Replace(staticItem.Name, "TORX", outputObj.Device.Type, -1)

		if len(staticItem.NextHop) != 0 {
			// Update null 0 static route
			tmp = append(tmp, StaticNetworkType{
				DstIPAddress: outputObj.getSupernetIPbyName(routeName),
				NextHop:      staticItem.NextHop,
				Name:         routeName,
			})
		} else if len(staticItem.DstIPAddress) != 0 {
			// update default route to border
			tmp = append(tmp, StaticNetworkType{
				DstIPAddress: outputObj.getSupernetIPbyName(staticItem.DstIPAddress),
				NextHop:      outputObj.getSwitchMgmtIPbyName(routeName),
				Name:         routeName,
			})
		}
	}
	r.Static.Network = tmp
}

func (r *RoutingType) updateStaticPrefixList(outputObj *OutputType) {
	for index, item := range r.Static.PrefixList {
		// Update Supernet Name to Supernet IP
		supernetIP := outputObj.getSupernetIPbyName(item.Supernet)
		r.Static.PrefixList[index].Supernet = supernetIP
	}
}

func (o *OutputType) getSupernetIPbyName(SupernetName string) string {
	for _, segment := range *o.Supernets {
		if segment.Name == SupernetName {
			return segment.Subnet
		}
	}
	return ""
}
