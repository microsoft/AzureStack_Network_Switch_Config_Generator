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
	if !o.Device.StaticRouting {
		// BGP Routing
		routingFrameworkObj.Bgp.BGPAsn = o.Device.Asn
		routingFrameworkObj.updateBgpNetwork(o)
		routerIDName := strings.Replace(routingFrameworkObj.Bgp.RouterID, "TORX", o.Device.Type, -1)
		routingFrameworkObj.Bgp.RouterID = o.searchSwitchMgmtIP(routerIDName)
		routingFrameworkObj.updateBgpNeighbor(o, inputJsonObj)
	}
	o.Routing = routingFrameworkObj
}

func (r *RoutingType) updateBgpNetwork(outputObj *OutputType) {
	for _, segment := range *outputObj.Supernets {
		for index, netname := range r.Bgp.IPv4Network {
			if segment.Name == netname {
				r.Bgp.IPv4Network[index] = segment.Subnet
			}
		}
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
		IPv4IPNet := outputObj.searchSwitchMgmtIP(nbrIPAddressName)
		r.Bgp.IPv4Neighbor[k].NeighborIPAddress = strings.Split(IPv4IPNet, "/")[0]
		r.Bgp.IPv4Neighbor[k].Description = nbrIPAddressName

		updateSourceName := replaceTORXName(v.UpdateSource, outputObj.Device.Type)
		r.Bgp.IPv4Neighbor[k].UpdateSource = outputObj.searchSwitchMgmtIP(updateSourceName)
	}

}
