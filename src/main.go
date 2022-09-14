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
	TOR1           = "TOR1"
	TOR2           = "TOR2"
	// No BMC Framework
	NOBMC = "nobmc"
	// Has BMC Framework
	HASBMC = "hasbmc"
	// PortType Name
	BMC_MGMT    = "BMCMgmt"
	IP          = "IP"
	ACCESS      = "Access"
	TRUNK       = "Trunk"
	INFRA_MGMT  = "InfraMgmt"
	SWITCH_MGMT = "SwitchMgmt"
	// Framework Module Name
	INTERFACE = "interface"
	ROUTING   = "routing"
	// File Extension
	JSON = "json"
	// Settings
	VPC               = "VPC"
	IBGP_PO           = "PO50"
	CHANNEL_GROUP     = "channel_group"
	CISCO_NATIVE_VLAN = "99"
)

// Logic: Input.json -> Object -Modify-> NewObject -> Output.json -> Template -> Config

func main() {
	// Set Log Output Options - 2022/08/24 21:51:10 main.go:58:
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/input_nobmc.json", "File path of switch deploy input.json")
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
