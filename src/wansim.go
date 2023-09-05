package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
)

func (o *OutputType) UpdateWANSIM(inputData InputData) {
	tmpReRouteNetwork := []string{}

	for _, rnetwork := range inputData.WANSIM.RerouteNetworks {
		_, vlanSubnetList := o.getSubnetByVlanGroupName(rnetwork)
		for _, subnet := range vlanSubnetList {
			fmt.Println(subnet)
			_, ipNet, err := net.ParseCIDR(subnet)
			if err != nil {
				fmt.Println(err)
				return
			}
			maskSize, _ := ipNet.Mask.Size()
			fmt.Println(GenSubnetsInNetwork(subnet, maskSize+1))
		}
		tmpReRouteNetwork = append(tmpReRouteNetwork, vlanSubnetList...)

	}

	o.WANSIM = inputData.WANSIM
	o.WANSIM.RerouteNetworks = tmpReRouteNetwork
}

// func DividSubnet(subnetStr string) {
// 	ip, ipNet, err := net.ParseCIDR(subnetStr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	maskSize, _ := ipNet.Mask.Size()
// 	newIPMask := net.CIDRMask(maskSize+1, 32)
// 	newNetMask := net.IP(newIPMask)

// 	firstSubnet := ip.Mask(newIPMask)
// 	secondSubnet := net.IP(make([]byte, len(firstSubnet)))
// 	copy(secondSubnet, firstSubnet)
// 	for i := len(secondSubnet) - 1; i >= 0; i-- {
// 		if secondSubnet[i] == 255 {
// 			secondSubnet[i] = 0
// 			continue
// 		}
// 		secondSubnet[i] += 255 - newNetMask[i]
// 		break
// 	}
// 	fmt.Printf("First subnet: %v/%d \n", firstSubnet, maskSize+1)
// 	fmt.Printf("Second subnet: %v/%d \n", secondSubnet, maskSize+1)
// }

func GenSubnetsInNetwork(netCIDR string, subnetMaskSize int) ([]string, error) {
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
