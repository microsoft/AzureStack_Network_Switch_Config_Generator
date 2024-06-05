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

	if strings.Contains(o.Switch.Make, "Dell") && strings.Contains(o.Switch.Type, "TOR") {
		o.updateDellPortGroup(interfaceJsonObj)
	}

	outputSwitchPorts := initSwitchPort(interfaceJsonObj, o.Switch.Make, o.Switch.Type, o.NodeCount)
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

func initSwitchPort(interfaceJsonObj *PortJson, switchMake, switchType string, nodeCount int) []PortType {
	outputSwitchPorts := []PortType{}
	portToIdx := map[string]int{}
	// the first forloop creates all the ports, calls them unused and shuts them down.
	for _, port := range interfaceJsonObj.Port {
		if strings.Contains(switchMake, "Dell") && strings.Contains(switchType, "TOR") {
			if strings.Contains(port.Mode, "10g-4x") {
				// 10g-4x port need to add ":1" while defining interface config
				outputSwitchPorts = append(outputSwitchPorts, PortType{
					Port:        port.Port + ":1",
					Idx:         port.Idx,
					Type:        port.Type,
					Shutdown:    true,
					Description: UNUSED,
					Function:    UNUSED,
					Mtu:         JUMBOMTU,
					UntagVlan:   UNUSED_VLANID,
					Mode:        port.Mode,
					PortGroup:   port.PortGroup,
				})
			} else {
				outputSwitchPorts = append(outputSwitchPorts, PortType{
					Port:        port.Port,
					Idx:         port.Idx,
					Type:        port.Type,
					Shutdown:    true,
					Description: UNUSED,
					Function:    UNUSED,
					Mtu:         JUMBOMTU,
					UntagVlan:   UNUSED_VLANID,
					Mode:        port.Mode,
					PortGroup:   port.PortGroup,
				})
			}
		} else {
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
		}
		portToIdx[port.Port] = port.Idx
	}

	// For function ports for compute set it to empty slice
	for i := range interfaceJsonObj.Function {
		isCompute := strings.EqualFold(interfaceJsonObj.Function[i].Function, "COMPUTE")
		isHostBMC := strings.EqualFold(interfaceJsonObj.Function[i].Function, "HOST_BMC")
		isStorage := strings.EqualFold(interfaceJsonObj.Function[i].Function, "STORAGE")
		if isCompute || isHostBMC || isStorage {
			// set interfaceJsonObj.Function[i].Port to an empty slice
			interfaceJsonObj.Function[i].Port = []string{}
		}
	}
	// function ports for compute set it to the number of nodes
	for i := range interfaceJsonObj.Function {
		isCompute := strings.EqualFold(interfaceJsonObj.Function[i].Function, "COMPUTE")
		isHostBMC := strings.EqualFold(interfaceJsonObj.Function[i].Function, "HOST_BMC")
		isStorage := strings.EqualFold(interfaceJsonObj.Function[i].Function, "STORAGE")
		if isCompute || isHostBMC {
			for key, value := range portToIdx {
				if value <= nodeCount {
					interfaceJsonObj.Function[i].Port = append(interfaceJsonObj.Function[i].Port, key)
					//fmt.Printf("Appending %s to Port of Function %d\n", key, i)
				}
			}
		}
		// For storage start at nodeCout to nodecount *2
		if isStorage {
			for key, value := range portToIdx {
				if value >= nodeCount+1 && value <= nodeCount*2 {
					interfaceJsonObj.Function[i].Port = append(interfaceJsonObj.Function[i].Port, key)
					//fmt.Printf("Appending %s to Port of Function %d\n", key, i)
				}
			}
		}
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
				if strings.HasSuffix(strings.ToUpper(vlanItem.VlanName), strings.ToUpper("_"+o.Switch.Type)) {
					STORAGE_VlanMap[vlanItem.VlanID] = vlanItem.VlanName
				}
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
			tmpPortObj.TagVlanList = COMPUTE_VlanList
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = COMPUTE
			tmpPortObj.Function = COMPUTE
		} else if strings.EqualFold(portItem.Function, COMPUTE) && strings.EqualFold(o.DeploymentPattern, HYPERCONVERGED) {
			// Switched Non Converged use both Compute and Storage Port Assignment
			tmpPortObj.UntagVlan = Compute_NativeVlanID
			tmpPortObj.TagVlanList = append(COMPUTE_VlanList, STORAGE_VlanList...)
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
			tmpPortObj.TagVlanList = COMPUTE_VlanList
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = fmt.Sprintf("%s-%s", o.DeploymentPattern, COMPUTE)
			tmpPortObj.Function = COMPUTE
		} else if strings.EqualFold(portItem.Function, STORAGE) && strings.EqualFold(o.DeploymentPattern, SWITCHED) {
			// Switched Non Converged use Storage Port Assignment
			tmpPortObj.UntagVlan = NATIVE_VLANID
			tmpPortObj.TagVlanList = STORAGE_VlanList
			tmpPortObj.Description = fmt.Sprintf("%s-%s", o.DeploymentPattern, STORAGE)
			tmpPortObj.Function = STORAGE
			tmpPortObj.Shutdown = false
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
			tmpPortObj.UntagVlan = NATIVE_VLANID
			tmpPortObj.TagVlanList = append(tmpPortObj.TagVlanList, BMC_VlanID)
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[TOR_BMC].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		} else if portItem.Function == TOR_BMC && len(o.SwitchBMC) == 0 {
			// No BMC
			tmpPortObj.UntagVlan = UNUSED_VLANID
			tmpPortObj.TagVlanList = nil
			tmpPortObj.Description = UNUSED
			tmpPortObj.Function = UNUSED
		} else if strings.EqualFold(portItem.Function, MLAG_PEER) {
			tmpPortObj.UntagVlan = NATIVE_VLANID
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[MLAG_PEER].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		}
		if len(tmpPortObj.TagVlanList) > 0 {
			tmpPortObj.TagVlanString = optimizeArray(tmpPortObj.TagVlanList)
		}
		o.Ports[i] = tmpPortObj
	}
}

func (o *OutputType) UpdateBMCSwitchPorts(VlanGroup map[string][]string) {
	for i, portItem := range o.Ports {
		tmpPortObj := portItem
		if strings.EqualFold(portItem.Function, HLHBMC) || strings.EqualFold(portItem.Function, HLHOS) || strings.EqualFold(portItem.Function, HOSTBMC) {
			tmpPortObj.UntagVlan = BMC_VlanID
			tmpPortObj.Shutdown = false
			tmpPortObj.Description = portItem.Function
			tmpPortObj.Function = portItem.Function
			tmpPortObj.Mtu = JUMBOMTU
		} else if strings.EqualFold(portItem.Function, TOR_BMC) {
			// Has BMC to TOR
			tmpPortObj.UntagVlan = NATIVE_VLANID
			tmpPortObj.TagVlanList = append(tmpPortObj.TagVlanList, BMC_VlanID)
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[TOR_BMC].PortChannelID,
			}
			tmpPortObj.Others = portOthers
			tmpPortObj.Shutdown = false
		}
		if len(tmpPortObj.TagVlanList) > 0 {
			tmpPortObj.TagVlanString = optimizeArray(tmpPortObj.TagVlanList)
		}
		o.Ports[i] = tmpPortObj
	}
}

func (o *OutputType) updateDellPortGroup(interfaceJsonObj *PortJson) {
	tmpPortGroup := []PortGroupType{}
	tmpPortGroupMap := map[string]PortGroupType{}
	for _, port := range interfaceJsonObj.Port {
		if len(port.PortGroup) > 0 {
			tmpPortGroupMap[port.PortGroup] = PortGroupType{
				PortGroup: port.PortGroup,
				Mode:      port.Mode,
				Type:      port.Type,
				Idx:       port.Idx,
			}
		}
	}
	for _, v := range tmpPortGroupMap {
		tmpPortGroup = append(tmpPortGroup, v)
	}
	sort.Slice(tmpPortGroup, func(i, j int) bool {
		return tmpPortGroup[i].Idx < tmpPortGroup[j].Idx
	})
	o.PortGroup = tmpPortGroup
}
