package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

var (
	ToolBuildVersion      = "1.2305.01"
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
	CRED_SCAN_PLACEHOLDER = "$CREDENTIAL_PLACEHOLDER$"

	JSONExtension          = ".json"
	YAMLExtension          = ".yaml"
	CONFIGExtension        = ".config"
	P2P_IBGP               = "P2P_IBGP"
	P2P_BORDER             = "P2P_BORDER"
	MLAG_PEER              = "MLAG_PEER"
	TOR_BMC                = "TOR_BMC"
	POID_P2P_IBGP          = "50"
	POID_MLAG_PEER         = "101"
	POID_TOR_BMC           = "102"
	UNUSED_VLANName        = "UNUSED_VLAN"
	UNUSED_VLANID          int
	NATIVE_VLANName        = "NativeVlan"
	NATIVE_VLANID          int
	BMC_VlanID             int
	Compute_NativeVlanName = "Infrastructure"
	Compute_NativeVlanID   int
	ANY                    = "Any"
	ANYNETWORK             = "0.0.0.0/0"
	WANSIM                 = "wansim"

	COMPUTE, STORAGE                         = "Compute", "Storage"
	SWITCHED, SWITCHLESS, HYPERCONVERGED     = "Switched", "Switchless", "HyperConverged"
	HLHBMC, HLHOS, HOSTBMC                   = "HLH_BMC", "HLH_OS", "HOST_BMC"
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
	wansimLibFolder := flag.String("wansimLibFolder", "../input/wansimLib", "Folder path of WAN-SIM solution templates")
	flag.StringVar(&Username, "username", "", "Username for switch configuration")
	flag.StringVar(&Password, "password", "", "Password for switch configuration")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputData := parseInputJson(*inputJsonFile)
	// Create device categrory map: Border, TOR, BMC, MUX based on Type
	DeviceTypeMap = inputData.createDeviceTypeMap()
	generateSwitchConfig(inputData, *switchLibFolder, *wansimLibFolder, *outputFolder, DeviceTypeMap)
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

func generateSwitchConfig(inputData InputData, switchLibFolder, wansimLibFolder, outputFolder string, DeviceTypeMap map[string][]SwitchType) {
	// TOR Switch
	if len(DeviceTypeMap[TOR]) > 0 {
		for _, torItem := range DeviceTypeMap[TOR] {
			torOutput := &OutputType{}
			// Function sequence matters, because the object construct phase by phase
			// Add Build Version
			torOutput.ToolBuildVersion = ToolBuildVersion
			torOutput.UpdateSwitch(torItem, TOR, DeviceTypeMap)
			// fmt.Printf("%#v\n%#v\n", torOutput, inputData)
			torOutput.UpdateVlanAndL3Intf(inputData)
			torOutput.UpdateGlobalSetting(inputData)
			torOutput.UpdateDHCPIps(inputData)
			templateFolder, frameworkFolder := torOutput.parseFrameworkPath(switchLibFolder)
			torOutput.UpdatePortChannel(inputData)
			torOutput.ParseSwitchPort(frameworkFolder)
			torOutput.ParseRouting(frameworkFolder, inputData)
			torOutput.UpdateWANSIM(inputData)
			debugYAMLOutput := path.Join(outputFolder,"debug_yaml")
			torOutput.writeToYaml(debugYAMLOutput)
			switchConfigOutput := path.Join(outputFolder,"switch_config")
			createFolder(switchConfigOutput)
			torOutput.parseCombineTemplate(templateFolder, switchConfigOutput, torItem.Hostname)
			wansimOutput := path.Join(outputFolder,"wansim_vm_config")
			createFolder(wansimOutput)
			torOutput.parseEachTemplate(wansimLibFolder, wansimOutput)
			onlyNewConfigOnlyOutput := path.Join(outputFolder,"wansim_switch_config")
			createFolder(onlyNewConfigOnlyOutput)
			torOutput.parseSelectedTemplate(templateFolder, onlyNewConfigOnlyOutput)
		}
	} else {
		log.Fatalln(NO_Valid_TOR_Switch)
	}
	// BMC Switch
	if len(DeviceTypeMap[BMC]) > 0 {
		for _, bmcItem := range DeviceTypeMap[BMC] {
			bmcOutput := &OutputType{}
			// Add Build Version
			bmcOutput.ToolBuildVersion = ToolBuildVersion
			bmcOutput.UpdateSwitch(bmcItem, BMC, DeviceTypeMap)
			bmcOutput.UpdateVlanAndL3Intf(inputData)
			bmcOutput.UpdateGlobalSetting(inputData)
			templateFolder, frameworkFolder := bmcOutput.parseFrameworkPath(switchLibFolder)
			bmcOutput.UpdatePortChannel(inputData)
			bmcOutput.ParseSwitchPort(frameworkFolder)
			bmcOutput.ParseRouting(frameworkFolder, inputData)
			debugYAMLOutput := path.Join(outputFolder,"debug_yaml")
			bmcOutput.writeToYaml(debugYAMLOutput)
			switchConfigOutput := path.Join(outputFolder,"switch_config")
			createFolder(switchConfigOutput)
			bmcOutput.parseCombineTemplate(templateFolder, switchConfigOutput, bmcItem.Hostname)
		}
	}
}
