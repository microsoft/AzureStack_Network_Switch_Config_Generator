package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	outputNetwork := inputObj.outputNetwork()
	randomUsername := "aszadmin-" + generateRandomString(5, 0, 0, 0)
	randomPassword := generateRandomString(16, 3, 3, 3)

	// Decode the Device section
	for _, deviceItem := range inputObj.Device {
		// Key GenerateDeviceConfig for further processing, otherwise will skip.
		if deviceItem.GenerateDeviceConfig {
			outputObj := NewOutputObj()
			// Determine the switch category based on Device info.
			frameworkPath, templatePath := deviceItem.validateSwitchFolder(*switchFolder)
			log.Println(deviceItem.Hostname, frameworkPath, templatePath)
			outputObj.Network = outputNetwork
			outputObj.Device = deviceItem
			outputObj.Device.Username = randomUsername
			outputObj.Device.Password = randomPassword

			// Dynamic updating output object based on switch framework.
			outputObj.updateOutputObj(frameworkPath, templatePath, inputObj)

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

func NewInputObj() *InputType {
	return &InputType{}
}

func NewOutputObj() *OutputType {
	return &OutputType{}
}

func newInterfaceFrameworkObj() *InterfaceFrameworkType {
	return &InterfaceFrameworkType{}
}

func parseInputJSON(inputJsonFile string) *InputType {
	inputJsonObj := NewInputObj()
	bytes, err := ioutil.ReadFile(inputJsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, inputJsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	return inputJsonObj
}

func parseOutputJSON(outputJsonFile string) *OutputType {
	outputObj := NewOutputObj()
	bytes, err := ioutil.ReadFile(outputJsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, outputObj)
	if err != nil {
		log.Fatalln(err)
	}
	return outputObj
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

func (o *OutputType) parseInterfaceObj(frameworkPath string) {
	interfaceFrameJson := fmt.Sprintf("%s/interface.json", frameworkPath)
	InterfaceFrameworkObj := parseInterfaceJSON(interfaceFrameJson)
	o.parseInBandPortFramework(InterfaceFrameworkObj)
	o.parseVlanObj(InterfaceFrameworkObj)
}

func parseInterfaceJSON(interfaceFrameJson string) *InterfaceFrameworkType {
	InterfaceFrameworkObj := newInterfaceFrameworkObj()
	bytes, err := ioutil.ReadFile(interfaceFrameJson)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, InterfaceFrameworkObj)
	if err != nil {
		log.Fatalln(err)
	}
	return InterfaceFrameworkObj
}

func (i *InputType) isRoutingBGP() bool {
	for _, v := range i.Device {
		if v.Asn <= 0 {
			return false
		}
	}
	return true
}

func (i *InputType) outputNetwork() *[]NetworkOutputType {
	NetworkBytes, err := json.Marshal(i.Network)
	if err != nil {
		log.Fatalln(err)
	}
	return parseNetworkSection(NetworkBytes)
}

func (o *OutputType) updateOutputObj(frameworkPath, templatePath string, inputObj *InputType) {
	o.parseInterfaceObj(frameworkPath)
	if inputObj.isRoutingBGP() {
		o.parseBGPFramework(frameworkPath, inputObj)
	}
}
