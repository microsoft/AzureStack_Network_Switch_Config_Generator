package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	inputFolder   = "../input"
	inputJsonFile = inputFolder + "/input.json"
	outputFolder  = "../output"
)

const (
	FRAMEWORK = "framework"
	TEMPLATE  = "template"
)

// Logic: Input.json -> Object -Modify-> NewObject -> Output.json -> Template -> Config

func main() {

	inputJsonObj := parseInputJSON(inputJsonFile)

	// Pass network content to Interface function
	NetworkBytes, err := json.Marshal(inputJsonObj.Network)
	if err != nil {
		log.Fatalln(err)
	}
	subnetIPList := parseInterfaceSection(NetworkBytes)

	for _, deviceItem := range inputJsonObj.Device {
		// fmt.Println(deviceItem)
		if deviceItem.GenerateDeviceConfig {
			outputJsonObj := NewOutputJsonObj()
			frameworkPath, templatePath := deviceItem.validateInputFolder()
			// fmt.Println(frameworkPath, templatePath)
			outputJsonObj.Network = subnetIPList
			outputJsonObj.Device = deviceItem
			outputJsonObj.parseInterfaceObj(frameworkPath)
			if inputJsonObj.isRoutingBGP() {
				outputJsonObj.parseBGPFramework(frameworkPath, inputJsonObj)
			}
			// Generate Output Object Json for Template to Consume
			outputJsonName := outputFolder + "/" + outputJsonObj.Device.Hostname + ".json"
			writeToJson(outputJsonName, outputJsonObj)
			// Parse Switch Template to Config
			outputConfigName := outputFolder + "/" + outputJsonObj.Device.Hostname + ".config"
			outputJsonObj.parseTemplate(templatePath, outputConfigName)
		}
	}
}

func NewInputJsonObj() *InputJsonType {
	return &InputJsonType{}
}

func NewOutputJsonObj() *OutputJsonType {
	return &OutputJsonType{}
}

func newInterfaceFrameworkObj() *InterfaceFrameworkType {
	return &InterfaceFrameworkType{}
}

func (o *OutputJsonType) parseInterfaceObj(frameworkPath string) {
	// interfaceFrameJson := fmt.Sprintf("%s/interface_%s.json", frameworkPath, strings.ToLower(o.Device.Type))
	interfaceFrameJson := fmt.Sprintf("%s/interface.json", frameworkPath)
	InterfaceFrameworkObj := parseInterfaceJSON(interfaceFrameJson)
	o.parseInBandPortFramework(InterfaceFrameworkObj)
	o.parseVlanFramework(InterfaceFrameworkObj)
}

func parseInputJSON(inputJsonFile string) *InputJsonType {
	inputJsonObj := NewInputJsonObj()
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

func (d *DeviceType) validateInputFolder() (frameworkPath, templatePath string) {
	frameworkPath = fmt.Sprintf("%s/%s/%s/%s/%s", inputFolder, d.Make, d.Model, d.Firmware, FRAMEWORK)
	frameworkPath = strings.ToLower(frameworkPath)
	_, err := os.Stat(frameworkPath)
	if err != nil {
		log.Println(err)
	}
	templatePath = fmt.Sprintf("%s/%s/%s/%s/%s", inputFolder, d.Make, d.Model, d.Firmware, TEMPLATE)
	templatePath = strings.ToLower(templatePath)

	_, err = os.Stat(templatePath)
	if err != nil {
		log.Println(err)
	}
	return frameworkPath, templatePath
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

func (i *InputJsonType) isRoutingBGP() bool {
	for _, v := range i.Device {
		if v.Asn <= 0 {
			return false
		}
	}
	return true
}
