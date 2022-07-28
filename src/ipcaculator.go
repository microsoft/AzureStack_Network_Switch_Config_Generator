package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"sort"
)

func parseSupernetSection(switchSubnetBytes []byte) *[]SupernetOutputType {
	switchSubnet := []SupernetInputType{}
	outputResult := []SupernetOutputType{}
	err := json.Unmarshal(switchSubnetBytes, &switchSubnet)
	if err != nil {
		log.Fatalln(err)
	}
	for _, inputSubnet := range switchSubnet {
		outputSubnet := newOutputSubnet()

		if len(inputSubnet.Subnet) != 0 {

			// Validate IP Size
			maxIPSize := getMaxIPSize(inputSubnet.Subnet)

			err := validateIPSize(maxIPSize, inputSubnet)
			if err != nil {
				log.Fatalln(err)
			}

			// Generate All IP List
			ipList := generateIPList(inputSubnet.Subnet)
			if err != nil {
				log.Fatalln(err)
			}

			// Subnet Segment [large -> small] and Assign IP from IP List
			outputSubnet.assignIPFromList(ipList, inputSubnet)
		} else {
			outputSubnet.updateNoSubnetObj(inputSubnet)
		}
		// fmt.Println(outputSubnet)
		outputResult = append(outputResult, *outputSubnet)
	}
	return &outputResult
}

func newOutputSubnet() *SupernetOutputType {
	return &SupernetOutputType{}
}

func (o *SupernetOutputType) assignIPFromList(ipList *[]string, inputSubnet SupernetInputType) {
	inputAssignment := inputSubnet.SubnetAssignment
	sort.SliceStable(inputAssignment, func(i, j int) bool {
		return inputAssignment[i].IPSize > inputAssignment[j].IPSize
	})

	pointer := 0
	tmpAssign := []IPAssignmentOutputItem{}
	for _, v := range inputAssignment {
		ipRange := (*ipList)[pointer : pointer+v.IPSize]

		if len(v.IPAssignment) != 0 {
			for _, p := range v.IPAssignment {
				fullName := fmt.Sprintf("%s/%s", v.Name, p.Name)
				network := fmt.Sprintf("%s/%d", ipRange[p.Position], v.Netmask)
				tmpAssign = append(tmpAssign, IPAssignmentOutputItem{
					Name:      fullName,
					IPAddress: network,
				})
			}
		}
		pointer += v.IPSize
	}
	o.VlanID = inputSubnet.VlanID
	o.Group = inputSubnet.Group
	o.Name = inputSubnet.Name
	o.Subnet = inputSubnet.Subnet
	o.Shutdown = inputSubnet.Shutdown
	o.IPAssignment = tmpAssign
}

func (o *SupernetOutputType) updateNoSubnetObj(inputSubnet SupernetInputType) {
	o.VlanID = inputSubnet.VlanID
	o.Group = inputSubnet.Group
	o.Name = inputSubnet.Name
	o.Subnet = inputSubnet.Subnet
	o.Shutdown = inputSubnet.Shutdown
	o.IPAssignment = []IPAssignmentOutputItem{}
}

func getMaxIPSize(ipnet string) int {
	_, ipv4Net, err := net.ParseCIDR(ipnet)
	if err != nil {
		log.Fatalln(err)
	}
	iMask, iBits := ipv4Net.Mask.Size()
	return int(math.Pow(2, float64(iBits-iMask)))
}

func validateIPSize(maxIPSize int, inputSubnet SupernetInputType) error {
	var actualIPSize int
	for k, s := range inputSubnet.SubnetAssignment {
		space := int(math.Pow(2, float64(32-s.Netmask)))
		// log.Println(s.Name, space)
		inputSubnet.SubnetAssignment[k].IPSize = space
		actualIPSize += space
	}

	if actualIPSize > maxIPSize {
		err := fmt.Errorf("[Error] Max assigned nework size %d, but asked for %d IP", maxIPSize, actualIPSize)
		return err
	}

	// log.Printf("[PASS] Max assigned nework size %d, and looking for %d IP\n", maxIPSize, actualIPSize)
	return nil
}

func generateIPList(ipnet string) *[]string {
	ipv4Addr, ipv4Net, err := net.ParseCIDR(ipnet)
	if err != nil {
		log.Fatalln(err)
	}
	var ips []string
	for ip := ipv4Addr.Mask(ipv4Net.Mask); ipv4Net.Contains(ip); ipIncrease(ip) {
		ips = append(ips, ip.String())
	}
	return &ips
}

func ipIncrease(ip net.IP) {
	// input is net.IP type which limites the value range
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
