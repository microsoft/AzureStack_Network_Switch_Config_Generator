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
	outputInterface := map[string]InterfaceType{}
	for _, port := range interfaceJsonObj.Port {
		portName := port.Port
		outputInterface[portName] = InterfaceType{
			Port: port.Port,
			Type: port.Type,
		}
	}
	// Initial Interface Object Map
	o.Interfaces = outputInterface
	// Config Interface with Functions
	for _, funcItem := range interfaceJsonObj.Function {
		for _, port := range funcItem.Port {
			if obj, ok := o.Interfaces[port]; ok {
				obj.Description = funcItem.Function
				o.Interfaces[port] = obj
			}
		}
	}
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
