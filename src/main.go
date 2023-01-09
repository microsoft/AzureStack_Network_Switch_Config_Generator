package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	TOR                 = "TOR"
	BMC                 = "BMC"
	BORDER              = "BORDER"
	MUX                 = "MUX"
	UPLINK              = "UPLINK"
	DOWNLINK            = "DOWNLINK"
	VIPGATEWAY          = "Gateway"
	UNUSED              = "Unused"
	TEMPLATE            = "template"
	INTERFACEJSON       = "interface.json"
	JUMBOMTU            = 9216
	DefaultMTU          = 1500
	NO_Valid_TOR_Switch = "NO Valid TOR Switch Founded"
	JSONExtension       = ".json"
	CONFIGExtension     = ".config"
	IBGP_PEER           = "PortChannel50"
	MLAG_PEER           = "PortChannel101"
	TOR_BMC             = "PortChannel102"
	Username, Password  string
)

func init() {
	// Set Log Output Options - example: 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/sample_input.json", "File path of switch deploy input.json")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")
	switchLibFolder := flag.String("switchLib", "../input/switchLib", "Folder path of switch frameworks and templates")
	flag.StringVar(&Username, "username", "", "Username for switch configuration")
	flag.StringVar(&Password, "password", "", "Password for switch configuration")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputObj := parseInputJson(*inputJsonFile)
	inputData := inputObj.InputData
	// Create random credential for switch config if no input values
	if Username == "" || Password == "" {
		Username = "aszadmin-" + generateRandomString(5, 0, 0, 0)
		Password = generateRandomString(16, 3, 3, 3)
	}

	// Create device categrory map: Border, TOR, BMC, MUX based on Type
	DeviceTypeMap := inputData.createDeviceTypeMap()
	// TOR Switch
	if len(DeviceTypeMap[TOR]) > 0 {
		for _, torItem := range DeviceTypeMap[TOR] {
			torOutput := &OutputType{}
			// Function sequence matters, because the object construct phase by phase
			torOutput.UpdateSwitch(torItem, TOR, DeviceTypeMap)
			torOutput.UpdateVlan(inputData)
			torOutput.UpdateGlobalSetting(inputData)
			templateFolder, frameworkFolder := torOutput.parseFrameworkPath(*switchLibFolder)
			torOutput.ParseSwitchInterface(frameworkFolder)
			// Output JSON File for Debug
			torOutput.writeToJson(*outputFolder)
			torOutput.parseTemplate(templateFolder, *outputFolder)
		}
	} else {
		log.Fatalln(NO_Valid_TOR_Switch)
	}
	// BMC Switch
	if len(DeviceTypeMap[BMC]) > 0 {
		for _, bmdItem := range DeviceTypeMap[BMC] {
			bmcOutput := &OutputType{}
			bmcOutput.UpdateSwitch(bmdItem, BMC, DeviceTypeMap)
			bmcOutput.UpdateVlan(inputData)
			bmcOutput.UpdateGlobalSetting(inputData)
			templateFolder, frameworkFolder := bmcOutput.parseFrameworkPath(*switchLibFolder)
			bmcOutput.ParseSwitchInterface(frameworkFolder)
			// Output JSON File for Debug
			bmcOutput.writeToJson(*outputFolder)
			bmcOutput.parseTemplate(templateFolder, *outputFolder)
		}
	}
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
