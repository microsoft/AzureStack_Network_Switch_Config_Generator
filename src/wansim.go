package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

var (
	PrefixList_DefaultRoute = "PL-DEFAULT"
	PrefixList_AllRoute     = "PL-ALL"
	RouteMap_Default_In     = "RM-DEFAULT-IN"
	RouteMap_Default_Out    = "RM-DEFAULT-OUT"
	RouteMap_NoRoute_IN     = "RM-NO-ROUTE-IN"
)

func (o *OutputType) UpdateWANSIM(inputData InputData) {
	tmpReRouteNetwork := []string{}

	for _, rnetwork := range inputData.WANSIM.RerouteNetworks {
		_, vlanSubnetList := o.getSubnetByVlanGroupName(rnetwork)
		for _, subnet := range vlanSubnetList {
			_, ipNet, err := net.ParseCIDR(subnet)
			if err != nil {
				fmt.Println(err)
				return
			}
			maskSize, _ := ipNet.Mask.Size()
			newSubnets, err := DividSubnetsByGivenMaskSize(subnet, maskSize+1)
			if err != nil {
				fmt.Println(err)
				return
			}
			tmpReRouteNetwork = append(tmpReRouteNetwork, newSubnets...)
		}
	}

	o.WANSIM = inputData.WANSIM
	o.UpdateWANSIMGRE(inputData)
	o.WANSIM.RerouteNetworks = tmpReRouteNetwork
	o.UpdateWANSIMBGP()
	o.UpdateWANSIMPingTestList()
}

func (o *OutputType) UpdateWANSIMGRE(inputData InputData) {
	// GRE Tunnel is established between WAN-SIM Loopback and TOR BMC VLAN IP
	// Local IP
	o.WANSIM.GRE1.TunnelSrcIP = o.WANSIM.Loopback.IP
	o.WANSIM.GRE2.TunnelSrcIP = o.WANSIM.Loopback.IP
	// Remote IP
	for _, supernetItem := range inputData.Supernets {
		if strings.ToUpper(supernetItem.GroupName) == "LOOPBACK0" {
			for _, assignItem := range supernetItem.IPv4.Assignment {
				if strings.Contains(strings.ToUpper(assignItem.Name), strings.ToUpper(o.WANSIM.GRE1.Name)) {
					o.WANSIM.GRE1.TunnelDstIP = assignItem.IP
				} else if strings.Contains(strings.ToUpper(assignItem.Name), strings.ToUpper(o.WANSIM.GRE2.Name)) {
					o.WANSIM.GRE2.TunnelDstIP = assignItem.IP
				}
			}
		}
	}
}

func (o *OutputType) UpdateWANSIMBGP() {
	// Update Staick Route Map for Uplink Nbr
	tmpIPv4Nbr := []IPv4NeighborType{}
	for _, nbrObj := range o.WANSIM.BGP.IPv4Nbr {
		tmpNbrObj := nbrObj
		tmpNbrObj.RouteMapIn = RouteMap_Default_In
		tmpIPv4Nbr = append(tmpIPv4Nbr, tmpNbrObj)
	}
	// Update GRE1 and GRE2 in BGP Section
	RemoteRackASN := o.Routing.BGP.BGPAsn
	// Add GRE1
	tmpIPv4Nbr = append(tmpIPv4Nbr, IPv4NeighborType{
		NeighborAsn:       RemoteRackASN,
		NeighborIPAddress: o.WANSIM.GRE1.RemoteIP,
		Description:       "To_TOR1",
		EbgpMultiHop:      8,
		UpdateSource:      "gre1",
		RouteMapIn:        RouteMap_NoRoute_IN,
		RouteMapOut:       RouteMap_Default_Out,
	})
	// Add GRE2
	tmpIPv4Nbr = append(tmpIPv4Nbr, IPv4NeighborType{
		NeighborAsn:       RemoteRackASN,
		NeighborIPAddress: o.WANSIM.GRE2.RemoteIP,
		Description:       "To_TOR2",
		EbgpMultiHop:      8,
		UpdateSource:      "gre2",
		RouteMapIn:        RouteMap_NoRoute_IN,
		RouteMapOut:       RouteMap_Default_Out,
	})
	o.WANSIM.BGP.IPv4Nbr = tmpIPv4Nbr
}

func (o *OutputType) UpdateWANSIMPingTestList() {
	DefaultPingList := []string{"\"microsoft.com\"", "\"azure.com\"", "\"msk8s.api.cdp.microsoft.com\""}
	for _, v := range o.Vlans {
		if v.GroupName == Compute_NativeVlanName {
			DefaultPingList = append(DefaultPingList, "\""+v.VIPAddress+"\"")
		} else if v.GroupName == BMC {
			DefaultPingList = append(DefaultPingList, "\""+v.VIPAddress+"\"")
		}
	}
	PingListStr := strings.Join(DefaultPingList, ",")
	o.WANSIM.PingTest = PingListStr
}

func DividSubnetsByGivenMaskSize(netCIDR string, subnetMaskSize int) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(netCIDR)
	if err != nil {
		return nil, err
	}
	if !ip.Equal(ipNet.IP) {
		return nil, errors.New("netCIDR is not a valid network address")
	}
	netMaskSize, _ := ipNet.Mask.Size()
	if netMaskSize > int(subnetMaskSize) {
		return nil, errors.New("subnetMaskSize must be greater or equal than netMaskSize")
	}

	totalSubnetsInNetwork := math.Pow(2, float64(subnetMaskSize)-float64(netMaskSize))
	totalHostsInSubnet := math.Pow(2, 32-float64(subnetMaskSize))
	subnetIntAddresses := make([]uint32, int(totalSubnetsInNetwork))
	// first subnet address is same as the network address
	subnetIntAddresses[0] = ip2int(ip.To4())
	for i := 1; i < int(totalSubnetsInNetwork); i++ {
		subnetIntAddresses[i] = subnetIntAddresses[i-1] + uint32(totalHostsInSubnet)
	}

	subnetCIDRs := make([]string, 0)
	for _, sia := range subnetIntAddresses {
		subnetCIDRs = append(
			subnetCIDRs,
			int2ip(sia).String()+"/"+strconv.Itoa(int(subnetMaskSize)),
		)
	}
	return subnetCIDRs, nil
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		panic("cannot convert IPv6 into uint32")
	}
	return binary.BigEndian.Uint32(ip)
}
func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
