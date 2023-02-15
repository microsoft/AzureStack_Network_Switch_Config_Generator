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
	outputInterface := []PortType{}
	portToIdx := map[string]int{}
	for _, port := range interfaceJsonObj.Port {
		outputInterface = append(outputInterface, PortType{
			Port:        port.Port,
			Idx:         port.Idx,
			Type:        port.Type,
			Description: UNUSED,
			Function:    UNUSED,
			Shutdown:    true,
			Mtu:         JUMBOMTU,
			UntagVlan:   UNUSED_VLANID,
		})
		portToIdx[port.Port] = port.Idx
	}
	// Initial Interface Object Map
	maxIdx := len(outputInterface)
	// Config Interface with Functions
	for _, funcItem := range interfaceJsonObj.Function {
		for _, port := range funcItem.Port {
			idxKey := portToIdx[port]
			if idxKey <= maxIdx {
				portItem := outputInterface[idxKey]
				portItem.Description = funcItem.Function
				portItem.Function = funcItem.Function
				portItem.Shutdown = false
				outputInterface[idxKey] = portItem
			} else {
				log.Fatalf("Port %s is not found in interface.json", port)
			}
		}
	}
	o.Ports = outputInterface
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

func (o *OutputType) UpdateSwitchPortByFunction() {
	STORAGE_VlanMap := map[int]string{}
	COMPUTE_VlanMap := map[int]string{}
	for _, vlanItem := range o.Vlans {
		for _, key := range STORAGEGroupName {
			if strings.Contains(vlanItem.GroupID, key) {
				STORAGE_VlanMap[vlanItem.VlanID] = vlanItem.GroupID
			}
		}
		for _, key := range COMPUTEGroupName {
			if strings.Contains(vlanItem.GroupID, key) {
				COMPUTE_VlanMap[vlanItem.VlanID] = vlanItem.GroupID
			}
		}
	}
	COMPUTE_VlanList := []int{}
	for COMPUTE_VlanID := range COMPUTE_VlanMap {
		COMPUTE_VlanList = append(COMPUTE_VlanList, COMPUTE_VlanID)
	}

	STORAGE_VlanList := []int{}
	for STORAGE_VlanID := range STORAGE_VlanMap {
		STORAGE_VlanList = append(STORAGE_VlanList, STORAGE_VlanID)
	}
	sort.Ints(COMPUTE_VlanList)
	sort.Ints(STORAGE_VlanList)

	for i, portItem := range o.Ports {
		if portItem.Function == "COMPUTE" {
			o.Ports[i].UntagVlan = Infra_VlanID
			o.Ports[i].TagVlans = COMPUTE_VlanList
		} else if portItem.Function == "STORAGE" {
			o.Ports[i].UntagVlan = Native_VLANID
			o.Ports[i].TagVlans = STORAGE_VlanList
		} else if strings.Contains(portItem.Function, "P2P_Border") {
			l3IntfName := fmt.Sprintf("%s_%s", portItem.Function, o.Switch.Type)
			portIpAddress := fmt.Sprintf("%s/%d", o.L3Interfaces[l3IntfName].IPAddress, o.L3Interfaces[l3IntfName].Cidr)
			o.Ports[i].IPAddress = portIpAddress
			o.Ports[i].UntagVlan = 0
		} else if portItem.Function == P2P_IBGP {
			o.Ports[i].UntagVlan = 0
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[P2P_IBGP].PortChannelID,
			}
			o.Ports[i].Others = portOthers
		} else if portItem.Function == TOR_BMC {
			o.Ports[i].UntagVlan = 0
			o.Ports[i].TagVlans = append(o.Ports[i].TagVlans, BMC_VlanID)
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[TOR_BMC].PortChannelID,
			}
			o.Ports[i].Others = portOthers
		} else if portItem.Function == MLAG_PEER {
			o.Ports[i].UntagVlan = Native_VLANID
			portOthers := map[string]string{
				"ChannelGroup": o.PortChannel[MLAG_PEER].PortChannelID,
			}
			o.Ports[i].Others = portOthers
		}
	}
}
