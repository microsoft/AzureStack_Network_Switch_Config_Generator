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

func (o *OutputType) parseBGPFramework(frameworkPath string, inputJsonObj *InputType) {
	routingFrameJson := fmt.Sprintf("%s/routing.json", frameworkPath)
	routingFrameworkObj := parseRoutingJSON(routingFrameJson)
	// Set unused section to null
	routingFrameworkObj.Router.Static = nil
	routingFrameworkObj.RouteMap = nil
	// Update template used variables
	routingFrameworkObj.Router.Bgp.BGPAsn = o.Device.Asn
	routingFrameworkObj.updateBgpNetwork(o)
	routerIDName := strings.Replace(routingFrameworkObj.Router.Bgp.RouterID, "TORX", o.Device.Type, -1)
	routingFrameworkObj.Router.Bgp.RouterID = o.searchSwitchMgmtIP(routerIDName)
	routingFrameworkObj.updateBgpNeighbor(o, inputJsonObj)
	o.Routing = routingFrameworkObj
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

func (r *RoutingType) updateBgpNetwork(outputObj *OutputType) {
	for _, segment := range *outputObj.Network {
		for index, netname := range r.Router.Bgp.IPv4Network {
			if segment.Name == netname {
				// fmt.Println(netname, segment.Subnet)
				r.Router.Bgp.IPv4Network[index] = segment.Subnet
			}
		}
	}
}

func (r *RoutingType) updateBgpNeighbor(outputObj *OutputType, inputJsonObj *InputType) {
	for k, v := range r.Router.Bgp.IPv4Neighbor {
		nbrAsn, err := inputJsonObj.getBgpASN(v.NeighborAsn)
		if err != nil {
			log.Fatalln(err)
		}
		r.Router.Bgp.IPv4Neighbor[k].NeighborAsn = nbrAsn
		nbrIPAddressName := strings.Replace(v.NeighborIPAddress, "TORX", outputObj.Device.Type, -1)
		r.Router.Bgp.IPv4Neighbor[k].NeighborIPAddress = outputObj.searchSwitchMgmtIP(nbrIPAddressName)
		r.Router.Bgp.IPv4Neighbor[k].Description = nbrIPAddressName

		updateSourceName := strings.Replace(v.UpdateSource, "TORX", outputObj.Device.Type, -1)
		r.Router.Bgp.IPv4Neighbor[k].UpdateSource = outputObj.searchSwitchMgmtIP(updateSourceName)
	}

}

func (i *InputType) getBgpASN(deviceName string) (string, error) {
	for _, v := range i.Device {
		if v.Hostname == deviceName {
			return fmt.Sprint(v.Asn), nil
		}
	}
	return "", fmt.Errorf("%s BGP ASN is invalid", deviceName)
}
