package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func (o *OutputType) ParseSwitchPort(frameworkFolder string) {
	interfaceJsonPath := fmt.Sprintf("%s/%s", frameworkFolder, INTERFACEJSON)
	interfaceJsonObj := parseInterfaceJson(interfaceJsonPath)
	outputSwitchPorts := initSwitchPort(interfaceJsonObj)
	o.Ports = outputSwitchPorts
	if strings.Contains(o.Switch.Type, TOR) {
		o.UpdateTORSwitchPorts(interfaceJsonObj.VlanGroup)
	} else if strings.Contains(o.Switch.Type, BMC) {
		o.UpdateBMCSwitchPorts(interfaceJsonObj.VlanGroup)
	}

}

func parseInterfaceJson(interfaceJsonPath string) *PortJson {
	interfaceJsonObj := &PortJson{}
	bytes, err := ioutil.ReadFile(interfaceJsonPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, interfaceJsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	return interfaceJsonObj
}

func initSwitchPort(interfaceJsonObj *PortJson) []PortType {
	outputSwitchPorts := []PortType{}
	portToIdx := map[string]int{}
	for _, port := range interfaceJsonObj.Port {
		outputSwitchPorts = append(outputSwitchPorts, PortType{
			Port:        port.Port,
			Idx:         port.Idx,
			Type:        port.Type,
			Shutdown:    true,
			Description: UNUSED,
			Function:    UNUSED,
			Mtu:         JUMBOMTU,
			UntagVlan:   UNUSED_VLANID,
		})
		portToIdx[port.Port] = port.Idx
	}
	// Initial Interface Object Map
	maxIdx := len(outputSwitchPorts)
	// Config Interface with Functions
	for _, funcItem := range interfaceJsonObj.Function {
		for _, port := range funcItem.Port {
			// port index starts at 1, so -1 to set index to 0
			idxKey := portToIdx[port] - 1
			if idxKey <= maxIdx {
				portItem := outputSwitchPorts[idxKey]
				portItem.Description = funcItem.Function
				portItem.Function = funcItem.Function
				outputSwitchPorts[idxKey] = portItem
			} else {
				log.Fatalf("Port %s is not found in interface.json", port)
			}
		}
	}
	return outputSwitchPorts
}

func (o *OutputType) UpdateTORSwitchPorts(VlanGroup map[string][]string) {
	// Get Storage and Compute VlanList
	STORAGE_VlanMap := map[int]string{}
	COMPUTE_VlanMap := map[int]string{}

	for _, vlanItem := range o.Vlans {
		for _, key := range VlanGroup[STORAGE] {
			if strings.Contains(strings.ToUpper(vlanItem.GroupName), strings.ToUpper(key)) {
				STORAGE_VlanMap[vlanItem.VlanID] = vlanItem.GroupName
			}
		}
		for _, key := range VlanGroup[COMPUTE] {
			if strings.Contains(strings.ToUpper(vlanItem.GroupName), strings.ToUpper(key)) {
				COMPUTE_VlanMap[vlanItem.VlanID] = vlanItem.GroupName
			}
		}
	}

	var STORAGE_VlanList, COMPUTE_VlanList []int
	for COMPUTE_VlanID := range COMPUTE_VlanMap {
		COMPUTE_VlanList = append(COMPUTE_VlanList, COMPUTE_VlanID)
	}

	for STORAGE_VlanID := range STORAGE_VlanMap {
		STORAGE_VlanList = append(STORAGE_VlanList, STORAGE_VlanID)
	}
	sort.Ints(COMPUTE_VlanList)
	sort.Ints(STORAGE_VlanList)

	for i, portItem := range o.Ports {
		tmpPortObj := portItem
		if strings.EqualFold(portItem.Function, COMPUTE) && strings.EqualFold(o.DeploymentPattern, SWITCHLESS) {
			// Switched Non Converged use both Compute and Storage Port Assignment
			tmpPortObj.UntagVlan = Compute_NativeVlanID
			tmpPortObj.TagVlans = COMPUTE_VlanList
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = COMPUTE
			tmpPortObj.Function = COMPUTE
		} else if strings.EqualFold(portItem.Function, COMPUTE) && strings.EqualFold(o.DeploymentPattern, HYPERCONVERGED) {
			// Switched Non Converged use both Compute and Storage Port Assignment
			tmpPortObj.UntagVlan = Compute_NativeVlanID
			tmpPortObj.TagVlans = append(COMPUTE_VlanList, STORAGE_VlanList...)
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = o.DeploymentPattern
			tmpPortObj.Function = o.DeploymentPattern
		} else if strings.EqualFold(portItem.Function, STORAGE) && (strings.EqualFold(o.DeploymentPattern, HYPERCONVERGED) || strings.EqualFold(o.DeploymentPattern, SWITCHLESS)) {
			// Remove the
			tmpPortObj.Shutdown = true
			tmpPortObj.Description = UNUSED
			tmpPortObj.Function = UNUSED
		} else if strings.EqualFold(portItem.Function, COMPUTE) && strings.EqualFold(o.DeploymentPattern, SWITCHED) {
			// Switched Non Converged use Compute Port Assignment
			tmpPortObj.UntagVlan = Compute_NativeVlanID
			tmpPortObj.TagVlans = COMPUTE_VlanList
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = fmt.Sprintf("%s-%s", o.DeploymentPattern, COMPUTE)
			tmpPortObj.Function = COMPUTE
		} else if strings.EqualFold(portItem.Function, STORAGE) && strings.EqualFold(o.DeploymentPattern, SWITCHED) {
			// Switched Non Converged use Storage Port Assignment
			if strings.EqualFold(o.Switch.Make, "Cisco") {
				// Cisco NXOS Storage Native Vlan is dummy vlan 99
				tmpPortObj.UntagVlan = CISCOMLAG_NATIVEVLANID
			} else if strings.EqualFold(o.Switch.Make, "DellEMC") {
				// DellEMC Storage Native Vlan is shutdonw and unused
				tmpPortObj.UntagVlan = UNUSED_VLANID
			}
			tmpPortObj.TagVlans = STORAGE_VlanList
			tmpPortObj.Description = fmt.Sprintf("%s-%s", o.DeploymentPattern, STORAGE)
			tmpPortObj.Function = STORAGE
		} else if strings.Contains(strings.ToUpper(portItem.Function), P2P_BORDER) {
			// Uplink to Border
			l3IntfName := strings.ToUpper(fmt.Sprintf("%s_%s", portItem.Function, o.Switch.Type))
			portIpAddress := fmt.Sprintf("%s/%d", o.L3Interfaces[l3IntfName].IPAddress, o.L3Interfaces[l3IntfName].Cidr)
			tmpPortObj.IPAddress = portIpAddress
			tmpPortObj.UntagVlan = 0
			tmpPortObj.Shutdown = false
		} else if portItem.Function == P2P_IBGP {
			tmpPortObj.UntagVlan = 0
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[P2P_IBGP].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		} else if portItem.Function == TOR_BMC && len(o.SwitchBMC) > 0 {
			// Has BMC
			tmpPortObj.UntagVlan = 0
			tmpPortObj.TagVlans = append(tmpPortObj.TagVlans, BMC_VlanID)
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[TOR_BMC].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		} else if portItem.Function == TOR_BMC && len(o.SwitchBMC) == 0 {
			// No BMC
			tmpPortObj.UntagVlan = UNUSED_VLANID
			tmpPortObj.TagVlans = nil
			tmpPortObj.Description = UNUSED
			tmpPortObj.Function = UNUSED
		} else if portItem.Function == MLAG_PEER {
			tmpPortObj.UntagVlan = CISCOMLAG_NATIVEVLANID
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[MLAG_PEER].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		}
		o.Ports[i] = tmpPortObj
	}
}

func (o *OutputType) UpdateBMCSwitchPorts(VlanGroup map[string][]string) {
	for i, portItem := range o.Ports {
		tmpPortObj := portItem
		if strings.EqualFold(portItem.Function, HLHBMC) || strings.EqualFold(portItem.Function, HLHOS) {
			tmpPortObj.UntagVlan = BMC_VlanID
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = portItem.Function
			tmpPortObj.Function = portItem.Function
		}
		o.Ports[i] = tmpPortObj
	}
}
