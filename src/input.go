package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func newInputObj() *InputType {
	return &InputType{}
}

func parseInputJSON(inputJsonFile string) *InputType {
	inputJsonObj := newInputObj()
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

func newInterfaceFrameworkObj() *InterfaceFrameworkType {
	return &InterfaceFrameworkType{}
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

func (i *InputType) parseSupernets() *[]SupernetOutputType {
	SupernetsBytes, err := json.Marshal(i.Supernets)
	if err != nil {
		log.Fatalln(err)
	}
	return parseSupernetSection(SupernetsBytes)
}

func (i *InputType) getBgpASN(typeName string) (string, error) {
	for _, v := range i.Devices {
		if v.Type == typeName {
			return fmt.Sprint(v.Asn), nil
		}
	}
	return "", fmt.Errorf("%s BGP ASN is invalid", typeName)
}
