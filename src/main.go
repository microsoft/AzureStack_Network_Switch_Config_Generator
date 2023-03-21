package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	TOR, BMC, BORDER, MUX = "TOR", "BMC", "BORDER", "MUX"
	UPLINK, DOWNLINK      = "UPLINK", "DOWNLINK"
	VIPGATEWAY            = "Gateway"
	UNUSED                = "Unused"
	TEMPLATE              = "template"
	INTERFACEJSON         = "interface.json"
	JUMBOMTU              = 9216
	DefaultMTU            = 1500
	BGP, STATIC           = "BGP", "STATIC"
	Username, Password    string

	JSONExtension            = ".json"
	CONFIGExtension          = ".config"
	P2P_IBGP                 = "P2P_IBGP"
	P2P_BORDER               = "P2P_BORDER"
	MLAG_PEER                = "MLAG_PEER"
	TOR_BMC                  = "TOR_BMC"
	POID_P2P_IBGP            = "50"
	POID_MLAG_PEER           = "101"
	POID_TOR_BMC             = "102"
	UNUSED_VLANName          = "UNUSED_VLAN"
	UNUSED_VLANID            int
	CISCOMLAG_NATIVEVLANNAME = "Cisco_MLAG_NativeVlan"
	CISCOMLAG_NATIVEVLANID   int
	BMC_VlanID               int
	Compute_NativeVlanName   = "Management"
	Compute_NativeVlanID     int
	ANY                      = "Any"
	ANYNETWORK               = "0.0.0.0/0"

	COMPUTE, STORAGE                         = "Compute", "Storage"
	SWITCHED, SWITCHLESS, HYPERCONVERGED     = "Switched", "Switchless", "Hyperconverged"
	HLHBMC, HLHOS, RESERVEDPDU               = "HLH_BMC", "HLH_OS", "Reserved_For_PDU"
	SWITCHUPLINK, SWITCHDOWNLINK, SWITCHPEER = "SwitchUplink", "SwitchDownlink", "SwitchPeer"
	BMC_DEFAULT_ROUTE                        = "GlobalDefaultRoute"
	DeviceTypeMap                            map[string][]SwitchType

	NO_Valid_TOR_Switch = "NO Valid TOR Switch Founded"
)

func init() {
	// Set Log Output Options - example: 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/cisco_sample_input1.json", "File path of switch deploy input.json")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")
	switchLibFolder := flag.String("switchLib", "../input/switchLib", "Folder path of switch frameworks and templates")
	flag.StringVar(&Username, "username", "", "Username for switch configuration")
	flag.StringVar(&Password, "password", "", "Password for switch configuration")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputData := parseInputJson(*inputJsonFile)
	// Create device categrory map: Border, TOR, BMC, MUX based on Type
	DeviceTypeMap = inputData.createDeviceTypeMap()
	generateSwitchConfig(inputData, *switchLibFolder, *outputFolder, DeviceTypeMap)
}

func (o *OutputType) parseFrameworkPath(switchLibFolder string) (string, string) {
	makeLow := strings.ToLower(o.Switch.Make)
	modelLow := strings.ToLower(o.Switch.Model)
	// Template Folder Path
	templateFolder := fmt.Sprintf("%s/%s/%s/%s", switchLibFolder, makeLow, o.Switch.Firmware, TEMPLATE)
	_, err := os.Stat(templateFolder)
	if err != nil {
		log.Println(err)
	}
	// Framework Folder Path
	frameworkFolder := fmt.Sprintf("%s/%s/%s/%s", switchLibFolder, makeLow, o.Switch.Firmware, modelLow)
	_, err = os.Stat(frameworkFolder)
	if err != nil {
		log.Println(err)
	}
	return templateFolder, frameworkFolder
}

func generateSwitchConfig(inputData InputData, switchLibFolder string, outputFolder string, DeviceTypeMap map[string][]SwitchType) {
	// TOR Switch
	if len(DeviceTypeMap[TOR]) > 0 {
		for _, torItem := range DeviceTypeMap[TOR] {
			torOutput := &OutputType{}
			// Function sequence matters, because the object construct phase by phase
			torOutput.UpdateSwitch(torItem, TOR, DeviceTypeMap)
			// fmt.Printf("%#v\n%#v\n", torOutput, inputData)
			torOutput.UpdateVlanAndL3Intf(inputData)
			torOutput.UpdatePortChannel(inputData)
			// fmt.Printf("%#v\n", torOutput)
			torOutput.UpdateGlobalSetting(inputData)
			templateFolder, frameworkFolder := torOutput.parseFrameworkPath(switchLibFolder)
			torOutput.ParseSwitchPort(frameworkFolder)
			torOutput.ParseRouting(frameworkFolder, inputData)
			// Output JSON File for Debug
			torOutput.writeToJson(outputFolder)
			torOutput.parseTemplate(templateFolder, outputFolder)
		}
	} else {
		log.Fatalln(NO_Valid_TOR_Switch)
	}
	// BMC Switch
	if len(DeviceTypeMap[BMC]) > 0 {
		for _, bmdItem := range DeviceTypeMap[BMC] {
			bmcOutput := &OutputType{}
			bmcOutput.UpdateSwitch(bmdItem, BMC, DeviceTypeMap)
			bmcOutput.UpdateVlanAndL3Intf(inputData)
			bmcOutput.UpdatePortChannel(inputData)
			bmcOutput.UpdateGlobalSetting(inputData)
			templateFolder, frameworkFolder := bmcOutput.parseFrameworkPath(switchLibFolder)
			bmcOutput.ParseSwitchPort(frameworkFolder)
			bmcOutput.ParseRouting(frameworkFolder, inputData)
			// Output JSON File for Debug
			bmcOutput.writeToJson(outputFolder)
			bmcOutput.parseTemplate(templateFolder, outputFolder)
		}
	}
}
