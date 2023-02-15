package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"
)

func (o *OutputType) writeToJson(outputFolder string) {
	// Create folder if not existing
	createFolder(outputFolder)
	jsonFile := outputFolder + "/" + o.Switch.Hostname + JSONExtension
	b, err := json.MarshalIndent(o, "", " ")
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

func (o *OutputType) parseTemplate(templateFolder, outputFolder string) {
	configFilePath := outputFolder + "/" + o.Switch.Hostname + CONFIGExtension
	t, err := template.ParseFiles(
		templateFolder+"/AllConfig.go.tmpl",
		templateFolder+"/hostname.go.tmpl",
		templateFolder+"/stig.go.tmpl",
		templateFolder+"/qos.go.tmpl",
		templateFolder+"/vlan.go.tmpl",
		templateFolder+"/portchannel.go.tmpl",
		templateFolder+"/port.go.tmpl",
		templateFolder+"/settings.go.tmpl",
		templateFolder+"/prefixlist.go.tmpl",
		templateFolder+"/bgp.go.tmpl",
	)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	err = t.Execute(f, o)
	if err != nil {
		log.Fatalln(err)
	}
	f.Close()
}
