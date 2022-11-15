package main

import (
	"flag"
	"log"
)

var (
	TOR                      = "TOR"
	BMC                      = "BMC"
	BORDER                   = "BORDER"
	MUX                      = "MUX"
	UPLINK                   = "UPLINK"
	DOWNLINK                 = "DOWNLINK"
	VIPGATEWAY               = "Gateway"
	UNUSED                   = "Unused"
	JUMBOMTU                 = 9216
	DefaultMTU               = 1500
	NO_Valid_TOR_Switch      = "NO Valid TOR Switch Founded"
	JSONExtension            = ".json"
	CONFIGExtension          = ".config"
	ranUsername, ranPassword string
)

func init() {
	// Set Log Output Options - example: 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/lab_input.json", "File path of switch deploy input.json")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")
	templateFolder := flag.String("switchFolder", "../input/switchfolder/cisco/template", "Folder path of switch frameworks and templates")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputObj := parseInputJson(*inputJsonFile)
	inputData := inputObj.InputData
	// Create random credential for switch config
	ranUsername = "aszadmin-" + generateRandomString(5, 0, 0, 0)
	ranPassword = generateRandomString(16, 3, 3, 3)

	// Create device categrory map: Border, TOR, BMC, MUX based on Type
	DeviceTypeMap := inputData.createDeviceTypeMap()
	// TOR Switch
	if len(DeviceTypeMap[TOR]) > 0 {
		for _, torItem := range DeviceTypeMap[TOR] {
			torOutput := &OutputType{}
			torOutput.UpdateSwitch(torItem, TOR, DeviceTypeMap)
			torOutput.UpdateGlobalSetting(inputData)
			torOutput.UpdateVlan(inputData)
			// Output JSON File for Debug
			torOutput.writeToJson(*outputFolder)
			torOutput.parseTemplate(*templateFolder, *outputFolder)
		}
	} else {
		log.Fatalln(NO_Valid_TOR_Switch)
	}
	// BMC Switch
	if len(DeviceTypeMap[BMC]) > 0 {
		for _, bmdItem := range DeviceTypeMap[BMC] {
			bmcOutput := &OutputType{}
			bmcOutput.UpdateSwitch(bmdItem, BMC, DeviceTypeMap)
			bmcOutput.UpdateGlobalSetting(inputData)
			bmcOutput.UpdateVlan(inputData)
			// Output JSON File for Debug
			bmcOutput.writeToJson(*outputFolder)
			bmcOutput.parseTemplate(*templateFolder, *outputFolder)
		}
	}
}
