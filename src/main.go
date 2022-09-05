package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	// Global
	FRAMEWORK = "framework"
	TEMPLATE  = "template"
	// Device Type
	DeviceType_TOR = "TOR"
	DeviceType_BMC = "BMC"
	// No BMC Framework
	NOBMC                         = "nobmc"
	TOR_NOBMC_INTERFACE_FRAMEWORK = "tor_nobmc_interface.json"
	TOR_NOBMC_ROUTING_FRAMEWORK   = "tor_nobmc_routing.json"
	// Has BMC Framework
	HASBMC                         = "hasbmc"
	TOR_HASBMC_INTERFACE_FRAMEWORK = "tor_hasbmc_interface.json"
	BMC_INTERFACE_FRAMEWORK        = "bmc_interface.json"
	// PortType Name
	PortType_BMC_MGMT   = "BMCMgmt"
	PortType_IP         = "IP"
	PortType_ACCESS     = "Access"
	PortType_TRUNK      = "Trunk"
	PortType_INFRA_MGMT = "InfraMgmt"
	// Framework Module Name
	INTERFACE = "interface"
	ROUTING   = "routing"
	// File Extension
	JSON = "json"
)

// Logic: Input.json -> Object -Modify-> NewObject -> Output.json -> Template -> Config

func main() {
	// Set Log Output Options - 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/input_hasbmc.json", "File path of switch deploy input.json")
	switchFolder := flag.String("switchFolder", "../input/switchfolder", "Folder path of switch frameworks and templates")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")

	// Covert input.json to Go Object
	inputObj := parseInputJSON(*inputJsonFile)

	// Decode Network section. (Pass the raw bytes instead of Obj, because trying to detach the ipcaculator function for future open source.)
	outputSupernets := inputObj.parseSupernets()
	// Create random credential for switch config
	randomUsername := "aszadmin-" + generateRandomString(5, 0, 0, 0)
	randomPassword := generateRandomString(16, 3, 3, 3)

	// Decode the Device section
	for _, deviceItem := range inputObj.Devices {
		// Key GenerateDeviceConfig for further processing, otherwise will skip.
		if deviceItem.GenerateDeviceConfig {
			outputObj := newOutputObj()
			// Determine the switch category based on Device info.
			frameworkPath := deviceItem.validateFrameworkPath(*switchFolder)
			templatePath := deviceItem.validateTemplatePath(*switchFolder)
			log.Println(deviceItem.Hostname, frameworkPath, templatePath)
			outputObj.Supernets = outputSupernets
			outputObj.Device = deviceItem
			outputObj.Device.Username = randomUsername
			outputObj.Device.Password = randomPassword
			outputObj.IsNoBMC = inputObj.IsNoBMC

			// Dynamic updating output object based on switch framework.
			outputObj.updateOutputObj(frameworkPath, inputObj)

			// Update External Section
			outputObj.updateSettings(inputObj)

			// Generate JSON Output for Debug
			createFolder(*outputFolder)
			outputJsonName := *outputFolder + "/" + outputObj.Device.Hostname + ".json"
			writeToJson(outputJsonName, outputObj)

			// Generate Generic Configuration Output for Deployment
			outputConfigName := *outputFolder + "/" + outputObj.Device.Hostname + ".config"
			outputObj.parseTemplate(templatePath, outputConfigName)
		}
	}
}

func (d *DeviceType) validateTemplatePath(switchFolder string) string {
	templatePath := fmt.Sprintf("%s/%s/%s", switchFolder, d.Make, TEMPLATE)
	templatePath = strings.ToLower(templatePath)

	_, err := os.Stat(templatePath)
	if err != nil {
		log.Println(err)
	}
	return templatePath
}

func (d *DeviceType) validateFrameworkPath(switchFolder string) string {
	frameworkPath := fmt.Sprintf("%s/%s/%s/%s/%s", switchFolder, d.Make, d.Model, d.Firmware, FRAMEWORK)
	frameworkPath = strings.ToLower(frameworkPath)

	_, err := os.Stat(frameworkPath)
	if err != nil {
		log.Println(err)
	}
	return frameworkPath
}