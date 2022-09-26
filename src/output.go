package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

func newOutputObj() *OutputType {
	return &OutputType{}
}

// Framework Selection based on Device Type
func (o *OutputType) updateOutputObj(frameworkPath string, inputObj *InputType) {
	if o.Device.Type == DeviceType_BMC {
		o.parseInterfaceObj(frameworkPath, DeviceType_BMC)
		o.parseRoutingFramework(frameworkPath, DeviceType_BMC, inputObj)
	} else {
		o.parseInterfaceObj(frameworkPath, DeviceType_TOR)
		o.parseRoutingFramework(frameworkPath, DeviceType_TOR, inputObj)
	}
}

func (o *OutputType) parseInterfaceObj(frameworkPath, deviceType string) {
	var interfaceFrameJson string
	if o.IsNoBMC {
		interfaceFrameJson = fmt.Sprintf("%s/%s_%s_%s.%s", frameworkPath, deviceType, INTERFACE, NOBMC, JSON)
	} else {
		interfaceFrameJson = fmt.Sprintf("%s/%s_%s_%s.%s", frameworkPath, deviceType, INTERFACE, HASBMC, JSON)
	}
	InterfaceFrameworkObj := parseInterfaceJSON(strings.ToLower(interfaceFrameJson))
	o.parseInBandPortFramework(InterfaceFrameworkObj)
	o.parseVlanObj(InterfaceFrameworkObj)
	o.parseLoopbackObj(InterfaceFrameworkObj)
	o.parsePortchannelObj(InterfaceFrameworkObj)
}

func (o *OutputType) updateSettings(inputObj *InputType) {
	if len(inputObj.Settings) == 0 {
		return
	}
	settingMap := inputObj.Settings
	// Add VPC srcIP and dstIP
	po50TOR1IPName := fmt.Sprintf("%s_%s", TOR1, IBGP_PO)
	po50TOR1IP := o.getIPbyName(po50TOR1IPName, SWITCH_MGMT)
	po50TOR2IPName := fmt.Sprintf("%s_%s", TOR2, IBGP_PO)
	po50TOR2IP := o.getIPbyName(po50TOR2IPName, SWITCH_MGMT)
	if o.Device.Type == TOR1 {
		settingMap[VPC] = []string{po50TOR1IP, po50TOR2IP}
	} else {
		settingMap[VPC] = []string{po50TOR2IP, po50TOR1IP}
	}
	// Add BMC_MGMT interface to External Section
	for _, v := range o.Vlan {
		if v.Group == BMC_MGMT {
			vlanIntf := fmt.Sprintf("Vlan%d", v.VlanID)
			settingMap[BMC_MGMT] = []string{vlanIntf}
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
	var t *template.Template
	var err error
	if o.Device.Type == DeviceType_BMC {
		t, err = template.ParseFiles(
			// BMC
			templatePath+"/bmcConfig.go.tmpl",
			templatePath+"/header.go.tmpl",
			templatePath+"/stig.go.tmpl",
			templatePath+"/port.go.tmpl",
			templatePath+"/vlan.go.tmpl",
			templatePath+"/default.go.tmpl",
			templatePath+"/stp.go.tmpl",
			templatePath+"/settings.go.tmpl",
			templatePath+"/qos.go.tmpl",
			templatePath+"/portchannel.go.tmpl",
			templatePath+"/bmcStatic.go.tmpl",
		)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// TOR
		t, err = template.ParseFiles(
			templatePath+"/torConfig.go.tmpl",
			templatePath+"/bmcConfig.go.tmpl",
			templatePath+"/header.go.tmpl",
			templatePath+"/stig.go.tmpl",
			templatePath+"/port.go.tmpl",
			templatePath+"/vlan.go.tmpl",
			templatePath+"/default.go.tmpl",
			templatePath+"/bgp.go.tmpl",
			templatePath+"/torStatic.go.tmpl",
			templatePath+"/stp.go.tmpl",
			templatePath+"/settings.go.tmpl",
			templatePath+"/qos.go.tmpl",
			templatePath+"/vpc.go.tmpl",
			templatePath+"/portchannel.go.tmpl",
		)
		if err != nil {
			log.Fatalln(err)
		}
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
