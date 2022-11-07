package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	TOR                 = "TOR"
	BMC                 = "BMC"
	BORDER              = "BORDER"
	MUX                 = "MUX"
	UPLINK              = "UPLINK"
	DOWNLINK            = "DOWNLINK"
	NO_Valid_TOR_Switch = "NO Valid TOR Switch Founded"
)

func init() {
	// Set Log Output Options - example: 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func main() {
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/lab_input.json", "File path of switch deploy input.json")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")
	flag.Parse()
	// Covert input.json to Go Object, structs are defined in model.go
	inputObj := parseInputJson(*inputJsonFile)
	inputData := inputObj.InputData
	// Create random credential for switch config
	// randomUsername := "aszadmin-" + generateRandomString(5, 0, 0, 0)
	// randomPassword := generateRandomString(16, 3, 3, 3)

	// Create device categrory map: Border, TOR, BMC, MUX based on Type
	DeviceTypeMap := inputData.createDeviceTypeMap()
	fmt.Println(DeviceTypeMap)
	// TOR Switch
	if len(DeviceTypeMap[TOR]) > 0 {
		for _, deviceTor := range DeviceTypeMap[TOR] {
			torOutput := &OutputType{}
			torOutput.Switch = deviceTor
			torOutput.SwitchUplink = DeviceTypeMap[UPLINK]
			torOutput.SwitchDownlink = DeviceTypeMap[DOWNLINK]
			torOutput.SwitchBMC = DeviceTypeMap[BMC]
			fmt.Println(torOutput)
			// Output JSON File for Debug
			torOutput.writeToJson(*outputFolder)
		}
	} else {
		log.Fatalln(NO_Valid_TOR_Switch)
	}
	// BMC Switch
	if len(DeviceTypeMap[BMC]) > 0 {
		for _, deviceBMC := range DeviceTypeMap[BMC] {
			bmcOutput := &OutputType{}
			bmcOutput.Switch = deviceBMC
			bmcOutput.SwitchUplink = DeviceTypeMap[TOR]
			fmt.Println(bmcOutput)
			bmcOutput.writeToJson(*outputFolder)
		}
	}
}
