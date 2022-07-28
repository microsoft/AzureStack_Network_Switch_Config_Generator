package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	FRAMEWORK = "framework"
	TEMPLATE  = "template"
)

// Logic: Input.json -> Object -Modify-> NewObject -> Output.json -> Template -> Config

func main() {

	// Input Variables
	inputJsonFile := flag.String("inputJsonFile", "../input/input.json", "File path of switch deploy input.json")
	switchFolder := flag.String("switchFolder", "../input/switchfolder", "Folder path of switch frameworks and templates")
	outputFolder := flag.String("outputFolder", "../output", "Folder path of switch configurations")

	// Covert input.json to Go Object
	inputObj := parseInputJSON(*inputJsonFile)

	// Decode Network section. (Pass the raw bytes instead of Obj, because trying to detach the ipcaculator function for future open source.)
	outputSupernets := inputObj.parseSupernets()
	randomUsername := "aszadmin-" + generateRandomString(5, 0, 0, 0)
	randomPassword := generateRandomString(16, 3, 3, 3)

	// Decode the Device section
	for _, deviceItem := range inputObj.Devices {
		// Key GenerateDeviceConfig for further processing, otherwise will skip.
		if deviceItem.GenerateDeviceConfig {
			outputObj := newOutputObj()
			// Determine the switch category based on Device info.
			frameworkPath, templatePath := deviceItem.validateSwitchFolder(*switchFolder)
			log.Println(deviceItem.Hostname, frameworkPath, templatePath)
			outputObj.Supernets = outputSupernets
			outputObj.Device = deviceItem
			outputObj.Device.Username = randomUsername
			outputObj.Device.Password = randomPassword

			// Dynamic updating output object based on switch framework.
			outputObj.updateOutputObj(frameworkPath, templatePath, inputObj)

			// Update External Section
			outputObj.updateSettings(inputObj)

			// Generate JSON Output for Debug
			createFolder(*outputFolder)
			outputJsonName := *outputFolder + "/" + outputObj.Device.Hostname + ".json"
			writeToJson(outputJsonName, outputObj)

			// Generate Configuration Output for Deployment
			outputConfigName := *outputFolder + "/" + outputObj.Device.Hostname + ".config"
			outputObj.parseTemplate(templatePath, outputConfigName)
		}
	}

}

func (d *DeviceType) validateSwitchFolder(switchFolder string) (frameworkPath, templatePath string) {
	frameworkPath = fmt.Sprintf("%s/%s/%s/%s/%s", switchFolder, d.Make, d.Model, d.Firmware, FRAMEWORK)
	frameworkPath = strings.ToLower(frameworkPath)
	_, err := os.Stat(frameworkPath)
	if err != nil {
		log.Println(err)
	}
	templatePath = fmt.Sprintf("%s/%s/%s/%s/%s", switchFolder, d.Make, d.Model, d.Firmware, TEMPLATE)
	templatePath = strings.ToLower(templatePath)

	_, err = os.Stat(templatePath)
	if err != nil {
		log.Println(err)
	}
	return frameworkPath, templatePath
}
