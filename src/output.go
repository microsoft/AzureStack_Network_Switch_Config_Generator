package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func newOutputObj() *OutputType {
	return &OutputType{}
}

func (o *OutputType) parseInterfaceObj(frameworkPath string) {
	interfaceFrameJson := fmt.Sprintf("%s/interface.json", frameworkPath)
	InterfaceFrameworkObj := parseInterfaceJSON(interfaceFrameJson)
	o.parseInBandPortFramework(InterfaceFrameworkObj)
	o.parseVlanObj(InterfaceFrameworkObj)
}

func (o *OutputType) updateOutputObj(frameworkPath, templatePath string, inputObj *InputType) {
	o.parseInterfaceObj(frameworkPath)
	o.parseRoutingFramework(frameworkPath, inputObj)
}

func (o *OutputType) updateSettings(inputObj *InputType) {
	if len(inputObj.Settings) == 0 {
		return
	}
	settingMap := inputObj.Settings
	// Add OOB interface to External Section
	for _, v := range o.Vlan {
		if v.Group == "OOB" {
			vlanIntf := fmt.Sprintf("Vlan%d", v.VlanID)
			settingMap["OOB"] = []string{vlanIntf}
		}
	}
	o.Settings = settingMap
}

func parseOutputJSON(outputJsonFile string) *OutputType {
	outputObj := newOutputObj()
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

func createFolder(folderPath string) {
	_, err := os.Stat(folderPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
}

func writeToJson(jsonFile string, outputResult interface{}) {
	b, err := json.MarshalIndent(outputResult, "", " ")
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(jsonFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = f.Write(b)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}

func (o *OutputType) parseTemplate(templatePath, outputConfigName string) {
	t, err := template.ParseFiles(
		templatePath+"/allConfig.go.tmpl",
		templatePath+"/header.go.tmpl",
		templatePath+"/stig.go.tmpl",
		templatePath+"/port.go.tmpl",
		templatePath+"/vlan.go.tmpl",
		templatePath+"/default.go.tmpl",
		templatePath+"/bgp.go.tmpl",
		templatePath+"/static.go.tmpl",
		templatePath+"/stp.go.tmpl",
		templatePath+"/settings.go.tmpl",
		templatePath+"/qos.go.tmpl",
	)
	if err != nil {
		log.Fatalln(err)
	}

	// err = t.Execute(os.Stdout, o)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	f, err := os.OpenFile(outputConfigName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln("create file: ", err)
		return
	}

	err = t.Execute(f, o)
	if err != nil {
		log.Fatalln("execute: ", err)
		return
	}
	f.Close()
}
