package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func (o *OutputType) ParseSwitchInterface(templateFolder string) {
	interfaceJsonPath := fmt.Sprintf("%s/%s", templateFolder, INTERFACEJSON)
	interfaceJsonObj := parseInterfaceJson(interfaceJsonPath)
	outputInterface := []map[string]InterfaceType{}
	for _, port := range interfaceJsonObj.Port {
		portName := port.Port
		tmpInterfaceMap := map[string]InterfaceType{}
		tmpInterfaceMap[portName] = InterfaceType{
			Port: port.Port,
			Type: port.Type,
		}
		outputInterface = append(outputInterface, tmpInterfaceMap)
	}
	// Initial Interface Object Map
	o.Interfaces = outputInterface
	// Config Interface with Functions
}

func parseInterfaceJson(interfaceJsonPath string) *InterfaceJson {
	interfaceJsonObj := &InterfaceJson{}
	bytes, err := ioutil.ReadFile(interfaceJsonPath)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(bytes, interfaceJsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	return interfaceJsonObj
}
